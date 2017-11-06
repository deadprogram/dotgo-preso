package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"strconv"
	"time"

	"gocv.io/x/gocv"
)

const (
	None = iota
	Laplacian
	Threshold
	Scharr
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\tdotgo [deviceid] [soundtrack]")
		return
	}

	deviceID, _ := strconv.Atoi(os.Args[1])
	soundtrack := os.Args[2]

	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("dotGo")
	window.SetWindowProperty(gocv.WindowPropertyFullscreen, gocv.WindowFullscreen)
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	imgFinal := gocv.NewMat()
	defer imgFinal.Close()

	imgFirst := gocv.NewMat()
	defer imgFirst.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	// start the soundtrack
	exec.Command("mpg123", soundtrack).Start()

	t := time.Now()
	start := t
	effect := 0
	algorithmName := "None"

	for {
		if time.Since(start) > 38*time.Second {
			os.Exit(0)
		}

		if ok := webcam.Read(img); !ok {
			fmt.Printf("Error cannot read device %d\n", deviceID)
			return
		}

		if img.Empty() {
			continue
		}

		// rotate through algorithms to use
		if time.Since(t) > 2*time.Second {
			t = time.Now()
			effect++
			if effect >= 4 {
				effect = None
			}
			img.CopyTo(imgFirst)
		}

		switch effect {
		case None:
			algorithmName = "None"
			img.CopyTo(imgFinal)
		case Laplacian:
			algorithmName = "Laplacian"
			gocv.Laplacian(img, imgFinal, gocv.MatTypeCV16S, 15, 1, 0, gocv.BorderDefault)
		case Scharr:
			algorithmName = "Scharr"
			gocv.Scharr(img, imgFinal, gocv.MatTypeCV16S, 1, 0, 30, 0, gocv.BorderDefault)
			gocv.Scharr(imgFinal, imgFinal, gocv.MatTypeCV16S, 0, 1, 30, 0, gocv.BorderDefault)
		case Threshold:
			algorithmName = "AbsDiff/Threshold"
			gocv.AbsDiff(imgFirst, img, imgDelta)
			gocv.Threshold(imgDelta, imgFinal, 25, 255, gocv.ThresholdBinary)
		}

		// show which algorithm in use
		size := gocv.GetTextSize(algorithmName, gocv.FontHersheyPlain, 1.5, 2)
		gocv.Rectangle(imgFinal, image.Rect(30, 30-size.Y-5, 30+size.X, 30+size.Y+5), color.RGBA{0, 0, 0, 0}, -1)
		gocv.PutText(imgFinal, algorithmName, image.Pt(30, 30),
			gocv.FontHersheyPlain, 1.5, color.RGBA{255, 255, 255, 0}, 2)

		window.IMShow(imgFinal)
		gocv.WaitKey(1)
	}
}
