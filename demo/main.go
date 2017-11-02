package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"gocv.io/x/gocv"
)

const (
	None = iota
	Laplacian
	Scharr
	Threshold
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

	for {
		if time.Since(start) > 35*time.Second {
			os.Exit(0)
		}

		if ok := webcam.Read(img); !ok {
			fmt.Printf("Error cannot read device %d\n", deviceID)
			return
		}

		if img.Empty() {
			continue
		}

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
			img.CopyTo(imgFinal)
		case Laplacian:
			gocv.Laplacian(img, imgFinal, gocv.MatTypeCV16S, 15, 1, 0, gocv.BorderDefault)
		case Scharr:
			gocv.Scharr(img, imgFinal, gocv.MatTypeCV16S, 1, 0, 30, 0, gocv.BorderDefault)
			gocv.Scharr(imgFinal, imgFinal, gocv.MatTypeCV16S, 0, 1, 30, 0, gocv.BorderDefault)
		case Threshold:
			gocv.AbsDiff(imgFirst, img, imgDelta)
			gocv.Threshold(imgDelta, imgFinal, 25, 255, gocv.ThresholdBinary)
		}

		window.IMShow(imgFinal)
		gocv.WaitKey(1)
	}
}
