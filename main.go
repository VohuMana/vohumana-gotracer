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
            Fuzziness: 0.0,
            Color: color.RGBA {
                R:1,
                G:1,
                B:255,
                A: 255 },
            IsDiffuse: false } }
    sphere2 := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: 3.0,
            Y: 0.5,
            Z: -5.0 },
        Radius: 1.0,
        Properties: raytracer.Material {
            Fuzziness: 0.2,
            Color: color.RGBA {
                R: 1,
                G: 255,
                B: 1,
                A: 255 },
            IsDiffuse: false } }
    sphere3 := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: -2.0,
            Y: 0.5,
            Z: -5.0 },
        Radius: 1.0,
        Properties: raytracer.Material {
            Fuzziness: 0.1,
            Color: color.RGBA {
                R: 255,
                G: 255,
                B: 255,
                A: 255 },
            IsDiffuse: false } }
    largeSphere := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: 0.,
            Y: -101.0,
            Z: 0.0 },
        Radius: 100.0,
        Properties: raytracer.Material {
            Color: color.RGBA {
                R: 128,
                G: 128,
                B: 128,
                A: 255 },
            IsDiffuse: true } }
            
    sphere.Properties.Attenuation = raytracer.AsVector3(sphere.Properties.Color)
    sphere2.Properties.Attenuation = raytracer.AsVector3(sphere2.Properties.Color)
    sphere3.Properties.Attenuation = raytracer.AsVector3(sphere3.Properties.Color)
    largeSphere.Properties.Attenuation = raytracer.AsVector3(largeSphere.Properties.Color)
    
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
    raytracer.Scene.AddObject("largeSphere", largeSphere)
    
    raytracer.MaxBounces = 16
    raytracer.MaxRaysPerBounce = 3
    
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