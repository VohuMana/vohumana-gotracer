package main

import
(
	"image"
    "image/color"
	"image/png"
	"log"
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
        Properties: raytracer.Material {
            Reflectiveness: 0.85,
            Color: color.RGBA {
                R:0,
                G:0,
                B:255,
                A: 255 } } }
    sphere2 := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: 0.5,
            Y: 5.0,
            Z: -10.0 },
        Radius: 1.0,
        Properties: raytracer.Material {
            Reflectiveness: 0.5,
            Color: color.RGBA {
                R: 0,
                G: 255,
                B: 0,
                A: 255 } } }
    largeSphere := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: 0.,
            Y: -101.0,
            Z: 0.0 },
        Radius: 100.0,
        Properties: raytracer.Material {
            Reflectiveness: 0.3,
            Color: color.RGBA {
                R: 255,
                G: 0,
                B: 0,
                A: 255 } } }
    raytracer.Scene.AddObject("sphere1", sphere)
    raytracer.Scene.AddObject("sphere2", sphere2)
    raytracer.Scene.AddObject("largeSphere", largeSphere)
    
    raytracer.MaxBounces = 16
    
    rayTracedFrame := image.NewRGBA(bounds)
    
    for y := 0; y < ySize; y++ {
        for x := 0; x < xSize; x++ {
            u := float32(x) / float32(xSize)
            v := float32(y) / float32(ySize)
            
            r := raytracer.Ray{
                Origin: camera.Origin,
                Direction: upperLeftImageCorner.Add(camera.ImagePlaneHorizontal.Scale(u)).Add(camera.ImagePlaneVertical.Scale(v)).UnitVector() }
            
            rayTracedFrame.Set(x, y, raytracer.ShootRay(r, raytracer.Scene, 0))
        }
    }
    
    outFile, err := os.Create("rayframe.png")
    checkError(err)
    defer outFile.Close()
    
    err = png.Encode(outFile, rayTracedFrame)
    checkError(err)
}