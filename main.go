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

func rayColor(r raytracer.Ray) color.RGBA {
    t := 0.5 * (r.Direction.Y + 1.0)
    // Lerp from blue to white
    c := raytracer.Vector3{ X: 0.5, Y: 0.7, Z: 1.0 }.Scale(1.0 - t).Add(raytracer.Vector3{X:1.0, Y: 1.0, Z: 1.0}.Scale(t)); 
    return color.RGBA{uint8(c.X * math.MaxUint8), uint8(c.Y * math.MaxUint8), uint8(c.Z * math.MaxUint8), 255}
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
    horiz := raytracer.Vector3{
        X: 4.0,
        Y: 0.0,
        Z: 0.0 }
    vert := raytracer.Vector3{
        X: 0.0,
        Y: 4.0,
        Z: 0.0 }
    origin := raytracer.Vector3{
        X: 0.0,
        Y: 0.0,
        Z: 0.0 }
    rayTracedFrame := image.NewRGBA(bounds)
    
    for y := 0; y < ySize; y++ {
        for x := 0; x < xSize; x++ {
            u := float32(x) / float32(xSize)
            v := float32(y) / float32(ySize)
            
            r := raytracer.Ray{
                Origin: origin,
                Direction: lowerLeftImageCorner.Add(horiz.Scale(u)).Add(vert.Scale(v)).UnitVector() }
            
            rayTracedFrame.Set(x, y, rayColor(r))
        }
    }
    
    outFile, err := os.Create("rayframe.png")
    checkError(err)
    defer outFile.Close()
    
    err = png.Encode(outFile, rayTracedFrame)
    checkError(err)
}