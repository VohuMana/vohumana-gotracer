package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/vohumana/vohumana-gotracer/RayTracer"
)

var communicationChannel chan bool

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	seedVal := time.Now().UTC().UnixNano()
	rand.Seed(seedVal)
	fmt.Printf("Generating scene with seed val: %v", seedVal)

	// shellIncrement := 10.0
	// sphereRadius := 20.0
	// numSpheresPerRow := 35.0
	// numShells := 3
	// incrementValue := 360.0 / numSpheresPerRow
	// incrementValue = float64(raytracer.ConvertDegreesToRadians(float32(incrementValue)))

	// for shell := 0; shell < numShells; shell++ {

	// 	theta := float64(0.0)
	// 	phi := float64(0.0)

	// 	for i := 0; i < int(numSpheresPerRow); i++ {
	// 		var pos raytracer.Vector3

	// 		switch shell % 3 {
	// 		case 0:
	// 			pos.X = float32(sphereRadius * math.Cos(theta))
	// 			pos.Y = float32(sphereRadius * math.Sin(phi))
	// 		case 1:
	// 			pos.X = float32(sphereRadius * math.Cos(theta))
	// 			pos.Z = float32(sphereRadius * math.Sin(phi))
	// 		case 2:
	// 			pos.Y = float32(sphereRadius * math.Cos(theta))
	// 			pos.Z = float32(sphereRadius * math.Sin(phi))
	// 		}
	// 		sphere := raytracer.Sphere{
	// 			Origin: pos,
	// 			Radius: 2.0}
	// 		color := color.RGBA{
	// 			R: uint8(rand.Intn(255)),
	// 			G: uint8(rand.Intn(255)),
	// 			B: uint8(rand.Intn(255)),
	// 			A: 255}

	// 		switch rand.Int() % 3 {
	// 		case 0:
	// 			sphere.Properties = raytracer.NewLambertian(color)
	// 		case 1:
	// 			sphere.Properties = raytracer.NewMetal(color, rand.Float32(), 1.0+(rand.Float32()*300), rand.Float32())
	// 		case 2:
	// 			sphere.Properties = raytracer.NewPhong(color, rand.Float32(), 1.0+(rand.Float32()*70))
	// 			// sphere.Properties = raytracer.Dielectric{
	// 			// 	Attenuation: raytracer.Vector3{
	// 			// 		X: 1.0,
	// 			// 		Y: 1.0,
	// 			// 		Z: 1.0},
	// 			// 	RefractiveIndex: }
	// 		}

	// 		theta += incrementValue
	// 		phi += incrementValue

	// 		raytracer.Scene.AddObject(strconv.Itoa((i+1)*(shell+1)), sphere)
	// 	}

	// 	sphereRadius += shellIncrement
	// }

	numSpheres := 150
	for i := 0; i < numSpheres; i++ {
		var pos raytracer.Vector3
		pos.X = (rand.Float32() * 100.0) - 50.0
		pos.Y = (rand.Float32() * 100.0) - 50.0
		pos.Z = (rand.Float32() * 100.0) - 50.0

		radius := 1.0 + (rand.Float32() * 4.0)

		sphere := raytracer.Sphere{
			Origin: pos,
			Radius: radius}
		color := color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255}

		switch rand.Int() % 3 {
		case 0:
			sphere.Properties = raytracer.NewLambertian(color)
		case 1:
			sphere.Properties = raytracer.NewMetal(color, rand.Float32(), 1.0+float32(int32(rand.Float32()*300)), rand.Float32())
		case 2:
			sphere.Properties = raytracer.NewPhong(color, rand.Float32(), 1.0+float32(int32(rand.Float32()*70)))
			// sphere.Properties = raytracer.Dielectric{
			// 	Attenuation: raytracer.Vector3{
			// 		X: 1.0,
			// 		Y: 1.0,
			// 		Z: 1.0},
			// 	RefractiveIndex: }
		}

		raytracer.Scene.AddObject(strconv.Itoa(i+1), sphere)
	}

	numLights := 5
	radius := 50.0
	phi := float64(0.0)
	theta := float64(0.0)
	incrementValue := 360.0 / float64(numLights)
	incrementValue = float64(raytracer.ConvertDegreesToRadians(float32(incrementValue)))
	for i := 0; i < numLights; i++ {
		color := color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255}
		// X := (rand.Float32() * 300.0) - 150.0
		// Y := (rand.Float32() * 300.0) - 150.0
		// Z := (rand.Float32() * 300.0) - 150.0
		X := float32(radius * math.Cos(theta))
		Z := float32(radius * math.Sin(phi))
		light := raytracer.NewPointLight(color, raytracer.NewVector3(X, 500.0, Z), rand.Float32()*2.0)
		raytracer.Scene.AddLight(strconv.Itoa(i+1), light)

		phi += incrementValue
		theta += incrementValue
	}

	raytracer.ExportScene("GeneratedScene.json", "GeneratedLights.json")

	// _ = raytracer.CreateCameraFromPos(
	//     raytracer.Vector3 {
	//         X: 0.0,
	//         Y: 0.0,
	//         Z: -1.0 },
	//     raytracer.Vector3 {
	//         X: 0.0,
	//         Y: 0.0,
	//         Z: 0.0 },
	//     raytracer.Vector3 {
	//         X: 0.0,
	//         Y: 1.0,
	//         Z: 0.0 },
	//     120,
	//     4.0 / 3.0)
	// raytracer.ExportCamera("camera.json")

}
