package main

import (
	// "flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/vohumana/vohumana-gotracer/raytracer"
)

var communicationChannel chan bool

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var configFilename string
	// var sceneFilename string
	var cameraFilename string

	// Get command line parameters
	// flag.StringVar(&configFilename, "config", "", "JSON filename describing how the ray tracer should render")
	// flag.StringVar(&sceneFilename, "scene", "", "JSON filename containing the scene to render")
	// flag.StringVar(&cameraFilename, "camera", "", "JSON filename containing the camera position and stats")
	// flag.Parse()

	// if (configFilename == "" || sceneFilename == "" || cameraFilename == "") {
	// 	flag.PrintDefaults()
	// 	return
	// }
	configFilename = "jsonfiles\\config4kPretty.json"
	cameraFilename = "SceneGenerator\\camera.json"
	raytracer.ImportConfig(configFilename)
	// raytracer.ImportScene(sceneFilename)
	raytracer.ImportCamera(cameraFilename)

	light := raytracer.NewPointLight(
		raytracer.Vector3{
			X: 1.0,
			Y: 1.0,
			Z: 1.0},
		raytracer.Vector3{
			X: -10,
			Y: 0,
			Z: -4},
		1.0)

	raytracer.Scene.AddLight("light", light)

	redPhong := raytracer.NewPhong(
		color.RGBA{
			R: 255},
		0.0,
		4)

	greenPhong := raytracer.NewPhong(
		color.RGBA{
			G: 255},
		0.1,
		4)

	metal := raytracer.NewMetal(
		color.RGBA{
			R: 255,
			G: 255,
			B: 255},
		0.9,
		4,
		0.0)

	fuzzyMetal := raytracer.NewMetal(
		color.RGBA{
			R: 255,
			B: 255},
		0.98,
		4,
		0.7)

	// diamond := raytracer.NewDielectric(0.05, 4, 2.4)
	lambertian := raytracer.NewLambertian(color.RGBA{B: 255})

	topSphere := raytracer.NewSphere(
		raytracer.Vector3{
			X: 0,
			Y: 3,
			Z: -5},
		1.5,
		lambertian)

	bottomSphere := raytracer.NewSphere(
		raytracer.Vector3{
			X: 0,
			Y: -3,
			Z: -5},
		1.5,
		redPhong)

	leftSphere := raytracer.NewSphere(
		raytracer.Vector3{
			X: -3,
			Y: 0,
			Z: -5},
		1.5,
		greenPhong)

	rightSphere := raytracer.NewSphere(
		raytracer.Vector3{
			X: 3,
			Y: 0,
			Z: -5},
		1.5,
		metal)

	middleSphere := raytracer.NewSphere(
		raytracer.Vector3{
			X: 0,
			Y: 0,
			Z: -5},
		0.25,
		fuzzyMetal)

	raytracer.Scene.AddObject("topSphere", topSphere)
	raytracer.Scene.AddObject("bottomSphere", bottomSphere)
	raytracer.Scene.AddObject("leftSphere", leftSphere)
	raytracer.Scene.AddObject("rightSphere", rightSphere)
	raytracer.Scene.AddObject("middleSphere", middleSphere)

	xSize := raytracer.Settings.WidthInPixels
	ySize := raytracer.Settings.HeightInPixels
	bounds := image.Rectangle{image.Point{0, 0}, image.Point{xSize, ySize}}

	rayTracedFrame := image.NewRGBA(bounds)
	communicationChannel = make(chan bool)

	startTime := time.Now()
	fmt.Printf("Beginning ray trace at resolution %v x %v\n", xSize, ySize)
	for y := 0; y < ySize; y++ {
		go RayTraceScanLine(rayTracedFrame, y, xSize, ySize)
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

	outFile, err := os.Create("rayframe.png")
	checkError(err)
	defer outFile.Close()

	err = png.Encode(outFile, rayTracedFrame)
	checkError(err)
}

// RayTraceScanLine will perform ray tracing for a single line of the image
func RayTraceScanLine(frame *image.RGBA, y, maxX, maxY int) {
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

	communicationChannel <- true
}
