// What it does:
//
// This example detects motion using a delta threshold from the first frame,
// and then finds contours to determine where the object is located.
//
// Based on Adrian Rosebrock code located at:
// http://www.pyimagesearch.com/2015/06/01/home-surveillance-and-motion-detection-with-the-raspberry-pi-python-and-opencv/
//
// How to run:
//
// 		go run ./track1/main.go 0
//

package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\ttracker1 [camera ID]")
		return
	}

	// parse args
	deviceID, _ := strconv.Atoi(os.Args[1])

	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("Capture Window")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	firstImg := gocv.NewMat()
	defer firstImg.Close()

	imgGrey := gocv.NewMat()
	defer imgGrey.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	fmt.Printf("Start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(img); !ok {
			fmt.Printf("Error cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		status := "Unoccupied"
		gocv.CvtColor(img, imgGrey, gocv.ColorBGRAToGray)
		gocv.GaussianBlur(imgGrey, imgGrey, image.Pt(21, 21), 0, 0, gocv.BorderDefault)

		if firstImg.Empty() {
			imgGrey.CopyTo(firstImg)
		}

		// calculate delta
		gocv.AbsDiff(firstImg, imgGrey, imgDelta)

		// use threshold
		gocv.Threshold(imgDelta, imgThresh, 25, 255, gocv.ThresholdBinary)

		// dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		gocv.Dilate(imgThresh, imgThresh, kernel)

		// find contours
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for _, c := range contours {
			area := gocv.ContourArea(c)
			if area < 5000 {
				continue
			}

			status = "Occupied"
			rect := gocv.BoundingRect(c)
			gocv.Rectangle(img, rect, color.RGBA{255, 0, 0, 0}, 2)
		}

		gocv.PutText(img, status, image.Pt(10, 20),
			gocv.FontHersheyPlain, 1.2, color.RGBA{255, 0, 0, 0}, 2)

		window.IMShow(img)
		gocv.WaitKey(1)
	}
}
