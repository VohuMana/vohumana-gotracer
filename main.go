package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
	"github.com/vohumana/vohumana-gotracer/raytracer"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Get command line parameters
	var (
		configFilename = flag.String("config", "", "JSON filename describing how the ray tracer should render")
		sceneFilename = flag.String("scene", "", "JSON filename containing the scene to render")
		lightsFilename = flag.String("light", "", "JSON filename containing the lights to render in the scene")
		cameraFilename = flag.String("camera", "", "JSON filename containing the camera position and stats")
		outFilename = flag.String("outfile", "rayframe.png", "Filename to output to, output will be a png")
		runPprof = flag.Bool("profile", false, "Specifiy this parameter if you want to run the profiler.  Pprof output will be saved to rayProfile.prof")
	)
	flag.Parse()

	if *configFilename == "" || *sceneFilename == "" || *cameraFilename == "" || *lightsFilename == "" {
		flag.PrintDefaults()
		return
	}
	
	if (*runPprof) {
		profileFile, err := os.Create("rayProfile")
		checkError(err)
		
		pprof.StartCPUProfile(profileFile)
		defer pprof.StopCPUProfile()
	}

	raytracer.ImportConfig(*configFilename)
	raytracer.ImportScene(*sceneFilename, *lightsFilename)
	raytracer.ImportCamera(*cameraFilename)

	xSize := raytracer.Settings.WidthInPixels
	ySize := raytracer.Settings.HeightInPixels
	bounds := image.Rectangle{image.Point{0, 0}, image.Point{xSize, ySize}}

	rayTracedFrame := image.NewRGBA(bounds)
	communicationChannel := make(chan bool)

	startTime := time.Now()
	fmt.Printf("Beginning ray trace at resolution %v x %v\n", xSize, ySize)
	for y := 0; y < ySize; y++ {
		go RayTraceScanLine(rayTracedFrame, y, xSize, ySize, &communicationChannel)
	}

	fmt.Println("All routines are running, now waiting")

	previousPercent := uint8(0)
	completedRoutines := 0
	for completedRoutines != ySize {
		percentComplete := uint8((float32(completedRoutines) / float32(ySize)) * 100.0)

		if previousPercent != percentComplete {
			fmt.Printf("%v%% Complete\n", percentComplete)
			previousPercent = percentComplete
		}

		success := <-communicationChannel
		if success {
			completedRoutines++
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Render duration was: %v s", elapsedTime.Seconds())

	outFile, err := os.Create(*outFilename)
	checkError(err)
	defer outFile.Close()

	err = png.Encode(outFile, rayTracedFrame)
	checkError(err)
}

// RayTraceScanLine will perform ray tracing for a single line of the image
func RayTraceScanLine(frame *image.RGBA, y, maxX, maxY int, channel *chan bool) {
	for x := 0; x < maxX; x++ {
		var red, green, blue float32

		for s := uint32(0); s < raytracer.Settings.MaxAntialiasRays; s++ {
			u := (float32(x) + rand.Float32()) / float32(maxX)
			v := (float32(y) + rand.Float32()) / float32(maxY)

			r := raytracer.Ray{
				Origin:    raytracer.GlobalCamera.Origin,
				Direction: raytracer.GlobalCamera.UpperLeftCorner.Add(raytracer.GlobalCamera.ImagePlaneHorizontal.Scale(u)).Add(raytracer.GlobalCamera.ImagePlaneVertical.Scale(v)).Subtract(raytracer.GlobalCamera.Origin).UnitVector()}

			color := raytracer.ShootRay(r, raytracer.Scene, 0)

			red += float32(color.R)
			green += float32(color.G)
			blue += float32(color.B)
		}

		red /= float32(raytracer.Settings.MaxAntialiasRays)
		green /= float32(raytracer.Settings.MaxAntialiasRays)
		blue /= float32(raytracer.Settings.MaxAntialiasRays)

		c := color.RGBA{
			R: uint8(red),
			G: uint8(green),
			B: uint8(blue),
			A: 255}

		// Render upside down because the image is upside down
		frame.Set(x, y, c)
	}

	channel <- true
}
