package main

import
(
	"image"
    "image/color"
	"image/png"
	"log"
    "math"
	"os"
)

func checkError(err error) {
	if (err != nil) {
		log.Fatal(err)
    }
}

func main() {
    xSize := 200
    ySize := 100
    bounds := image.Rectangle{image.Point{0,0}, image.Point{xSize, ySize}}
    rayTracedFrame := image.NewRGBA(bounds)
    
    for y := 0; y < ySize; y++ {
        for x := 0; x < xSize; x++ {
            var r float32
            var g float32
            var b float32
            
            r = float32(x) / float32(xSize)
            g = float32(y) / float32(ySize)
            b = 0.2
            
            rayTracedFrame.Set(x, y, color.RGBA{uint8(r * math.MaxUint8), uint8(g * math.MaxUint8), uint8(b * math.MaxUint8), 255})
        }
    }
    
    outFile, err := os.Create("rayframe.png")
    checkError(err)
    defer outFile.Close()
    
    err = png.Encode(outFile, rayTracedFrame)
    checkError(err)
}