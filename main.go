package main

import
(
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

var globalCamera raytracer.Camera
var communicationChannel chan bool

func checkError(err error) {
	if (err != nil) {
		log.Fatal(err)
    }
}

func main() {
    raytracer.ImportConfig("config.json")
    raytracer.ImportScene("scene.json")
    xSize := raytracer.Settings.WidthInPixels
    ySize := raytracer.Settings.HeightInPixels
    bounds := image.Rectangle{image.Point{0,0}, image.Point{xSize, ySize}}
    globalCamera = raytracer.CreateCameraFromPos(
        raytracer.Vector3 {
            X: 0.5,
            Y: 0.5,
            Z: -5.0 },
        raytracer.Vector3 {
            X: -0.5,
            Y: 2.0,
            Z: -7.0 },
        raytracer.Vector3 {
            X: 0.0,
            Y: 1.0,
            Z: 0.0 },
        120, 
        float32(xSize) / float32(ySize))
    // globalCamera = raytracer.CreateCamera(120, float32(xSize) / float32(ySize))
    
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
                Origin: globalCamera.Origin,
                Direction: globalCamera.UpperLeftCorner.Add(globalCamera.ImagePlaneHorizontal.Scale(u)).Add(globalCamera.ImagePlaneVertical.Scale(v)).Subtract(globalCamera.Origin).UnitVector() }
                
            color := raytracer.ShootRay(r, raytracer.Scene, 0)
            
            red += float32(color.R)
            green += float32(color.G)
            blue += float32(color.B)
        }
        
        red /= float32(raytracer.Settings.MaxAntialiasRays)
        green /= float32(raytracer.Settings.MaxAntialiasRays)
        blue /= float32(raytracer.Settings.MaxAntialiasRays)
        
        c := color.RGBA {
            R: uint8(red),
            G: uint8(green),
            B: uint8(blue),
            A: 255 }
        
        // Render upside down because the image is upside down
        frame.Set(x, y, c)
    }
    
    communicationChannel <- true
}