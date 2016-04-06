package main

import
(
	"image"
    "image/color"
	"image/png"
	"log"
    "math"
	"os"
    "github.com/vohumana/vohumana-gotracer/raytracer"
)

func rayColor(r raytracer.Ray, o raytracer.CollidableObject) color.RGBA {
    collided, record := o.TestIntersection(r)
    if (collided) {
        c := record.Normal.Scale(0.5)
        c = c.Add(raytracer.Vector3{X:1.0, Y:1.0, Z:1.0})
        
        c.X = float32(math.Min(1.0, float64(c.X)))
        c.Y = float32(math.Min(1.0, float64(c.Y)))
        c.Z = float32(math.Min(1.0, float64(c.Z)))
        return c.AsColor()
    }
    
    t := 0.5 * (r.Direction.Y + 1.0)
    // Lerp from blue to white
    c := raytracer.Vector3{ X: 0.5, Y: 0.7, Z: 1.0 }.Scale(1.0 - t).Add(raytracer.Vector3{X:1.0, Y: 1.0, Z: 1.0}.Scale(t)); 
    return c.AsColor()
}

func checkError(err error) {
	if (err != nil) {
		log.Fatal(err)
    }
}

func main() {
    xSize := 3840
    ySize := 2160
    bounds := image.Rectangle{image.Point{0,0}, image.Point{xSize, ySize}}
    lowerLeftImageCorner := raytracer.Vector3{
        X: -2.0,
        Y: -1.0,
        Z: -1.0 }
    camera := raytracer.Camera{
        Origin: raytracer.Vector3{
            X: 0.0,
            Y: 0.0,
            Z: 0.0 },
        ImagePlaneHorizontal: raytracer.Vector3{
            X: 3.8,
            Y: 0.0,
            Z: 0.0 },
        ImagePlaneVertical: raytracer.Vector3{
            X: 0.0,
            Y: 2.1,
            Z: 0.0 } }
    sphere := raytracer.Sphere{
        Origin: raytracer.Vector3{
            X: 0.5,
            Y: 0.5,
            Z: -5.0 },
        Radius: 1.0 }
    rayTracedFrame := image.NewRGBA(bounds)
    
    for y := 0; y < ySize; y++ {
        for x := 0; x < xSize; x++ {
            u := float32(x) / float32(xSize)
            v := float32(y) / float32(ySize)
            
            r := raytracer.Ray{
                Origin: camera.Origin,
                Direction: lowerLeftImageCorner.Add(camera.ImagePlaneHorizontal.Scale(u)).Add(camera.ImagePlaneVertical.Scale(v)).UnitVector() }
            
            rayTracedFrame.Set(x, y, rayColor(r, sphere))
        }
    }
    
    outFile, err := os.Create("rayframe.png")
    checkError(err)
    defer outFile.Close()
    
    err = png.Encode(outFile, rayTracedFrame)
    checkError(err)
}