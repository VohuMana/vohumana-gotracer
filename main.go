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
    "github.com/vohumana/vohumana-gotracer/raytracer"
)

var globalCamera raytracer.Camera
var globalUpperLeftCorner raytracer.Vector3
var communicationChannel chan bool

func checkError(err error) {
	if (err != nil) {
		log.Fatal(err)
    }
}

func main() {
    xSize := 3840
    ySize := 2160
    bounds := image.Rectangle{image.Point{0,0}, image.Point{xSize, ySize}}
    blue := color.RGBA {
                R:1,
                G:1,
                B:255,
                A: 255 }
    green := color.RGBA {
                R: 1,
                G: 255,
                B: 1,
                A: 255 }
    white := color.RGBA {
                R: 255,
                G: 255,
                B: 255,
                A: 255 }
    grey := color.RGBA {
                R: 128,
                G: 128,
                B: 128,
                A: 255 }
    globalUpperLeftCorner = raytracer.Vector3{
        X: 2.0,
        Y: 1.0,
        Z: -1.0 }
    globalCamera = raytracer.Camera{
        Origin: raytracer.Vector3{
            X: 0.0,
            Y: 0.0,
            Z: 0.0 },
        ImagePlaneHorizontal: raytracer.Vector3{
            X: -3.84,
            Y: 0.0,
            Z: 0.0 },
        ImagePlaneVertical: raytracer.Vector3{
            X: 0.0,
            Y: -2.16,
            Z: 0.0 } }
    sphere := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: 0.5,
            Y: 0.5,
            Z: -5.0 },
        Radius: 1.0,
        Properties: raytracer.Metal {
            Fuzziness: 0.0,
            Color: blue,
            Attenuation: raytracer.AsVector3(blue) } }
    sphere2 := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: 3.0,
            Y: 0.5,
            Z: -5.0 },
        Radius: 1.0,
        Properties: raytracer.Metal {
            Fuzziness: 0.2,
            Color: green,
            Attenuation: raytracer.AsVector3(green) } }
    sphere3 := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: -2.0,
            Y: 0.5,
            Z: -5.0 },
        Radius: 1.0,
        Properties: raytracer.Metal {
            Fuzziness: 0.1,
            Color: white,
            Attenuation: raytracer.AsVector3(white) } }
    diamondSphere := raytracer.Sphere {
        Origin: raytracer.Vector3 {
            X: 0.0,
            Y: 0.0,
            Z: -2.0 },
        Radius: 0.25,
        Properties: raytracer.Dielectric {
            RefractiveIndex: 2.4,
            Attenuation: raytracer.AsVector3(white) } }
    largeSphere := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: 0.,
            Y: -101.0,
            Z: 0.0 },
        Radius: 100.0,
        Properties: raytracer.Lambertian {
            Color: grey,
            Attenuation: raytracer.AsVector3(grey) } }
    
    raytracer.SkyColorBottom = raytracer.AsVector3(color.RGBA {
        R: 255,
        G: 239,
        B: 138 })
    
    raytracer.SkyColorTop = raytracer.AsVector3(color.RGBA {
        R: 40,
        G: 105,
        B: 209 })
    
    raytracer.Scene.AddObject("sphere1", sphere)
    raytracer.Scene.AddObject("sphere2", sphere2)
    raytracer.Scene.AddObject("sphere3", sphere3)
    raytracer.Scene.AddObject("diamondSphere", diamondSphere)
    raytracer.Scene.AddObject("largeSphere", largeSphere)
    
    raytracer.MaxBounces = 5
    raytracer.MaxRaysPerBounce = 3
    raytracer.MaxAntialiasRays = 3
    
    rayTracedFrame := image.NewRGBA(bounds)
    communicationChannel = make(chan bool)
    
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
    
    outFile, err := os.Create("rayframe.png")
    checkError(err)
    defer outFile.Close()
    
    err = png.Encode(outFile, rayTracedFrame)
    checkError(err)
}

func RayTraceScanLine(frame *image.RGBA, y, maxX, maxY int) {
    for x := 0; x < maxX; x++ { 
        var red, green, blue float32
        
        for s := uint32(0); s < raytracer.MaxAntialiasRays; s++ {
            u := (float32(x) + rand.Float32()) / float32(maxX)
            v := (float32(y) + rand.Float32()) / float32(maxY)
                    
            r := raytracer.Ray{
                Origin: globalCamera.Origin,
                Direction: globalUpperLeftCorner.Add(globalCamera.ImagePlaneHorizontal.Scale(u)).Add(globalCamera.ImagePlaneVertical.Scale(v)).UnitVector() }
                
        color := raytracer.ShootRay(r, raytracer.Scene, 0)
        
        red += float32(color.R)
        green += float32(color.G)
        blue += float32(color.B)
        }
        
        red /= float32(raytracer.MaxAntialiasRays)
        green /= float32(raytracer.MaxAntialiasRays)
        blue /= float32(raytracer.MaxAntialiasRays)
        
        c := color.RGBA {
            R: uint8(red),
            G: uint8(green),
            B: uint8(blue),
            A: 255 }
        
        frame.Set(x, y, c)
    }
    
    communicationChannel <- true
}