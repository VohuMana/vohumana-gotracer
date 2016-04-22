package main

import
(
	"image"
    "image/color"
	"image/png"
	"log"
    "math/rand"
	"os"
    "github.com/vohumana/vohumana-gotracer/raytracer"
)

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
    upperLeftImageCorner := raytracer.Vector3{
        X: 2.0,
        Y: 1.0,
        Z: -1.0 }
    camera := raytracer.Camera{
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
            Y: 0.5,
            Z: -2.0 },
        Radius: 0.5,
        Properties: raytracer.Dielectric {
            RefractiveIndex: 1.3,
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
    
    for y := 0; y < ySize; y++ {
        for x := 0; x < xSize; x++ { 
           TraceRayForPoint(rayTracedFrame, x, y, xSize, ySize, camera, upperLeftImageCorner)
        }
    }
    
    outFile, err := os.Create("rayframe.png")
    checkError(err)
    defer outFile.Close()
    
    err = png.Encode(outFile, rayTracedFrame)
    checkError(err)
}

func TraceRayForPoint(frame *image.RGBA, x, y , maxX, maxY int, camera raytracer.Camera, upperLeftImageCorner raytracer.Vector3) {
    var red, green, blue float32
    
    for s := uint32(0); s < raytracer.MaxAntialiasRays; s++ {
        u := (float32(x) + rand.Float32()) / float32(maxX)
        v := (float32(y) + rand.Float32()) / float32(maxY)
                
        r := raytracer.Ray{
            Origin: camera.Origin,
            Direction: upperLeftImageCorner.Add(camera.ImagePlaneHorizontal.Scale(u)).Add(camera.ImagePlaneVertical.Scale(v)).UnitVector() }
            
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