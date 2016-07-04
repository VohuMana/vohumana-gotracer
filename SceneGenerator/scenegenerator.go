package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
	"github.com/vohumana/vohumana-gotracer/RayTracer"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func spawnLightsInXZPlane(Y float32, numLights uint) {
	radius := 50.0
	phi := float64(0.0)
	theta := float64(0.0)
	incrementValue := 360.0 / float64(numLights)
	incrementValue = float64(raytracer.ConvertDegreesToRadians(float32(incrementValue)))
	for i := uint(0); i < numLights; i++ {
		color := getRandomColor()

		X := float32(radius * math.Cos(theta))
		Z := float32(radius * math.Sin(phi))
		
		light := raytracer.NewPointLight(color, raytracer.NewVector3(X, Y, Z), rand.Float32() * 100.0)
		raytracer.Scene.AddLight(strconv.Itoa(int(i+1)), light)

		phi += incrementValue
		theta += incrementValue
	}
}

func spawnSpheresInCube(height, width, depth float32, numSpheres uint) {
	for i := uint(0); i < numSpheres; i++ {
		var pos raytracer.Vector3
		pos.X = (rand.Float32() * width) - (width / 2.0)
		pos.Y = (rand.Float32() * height) - (height / 2.0)
		pos.Z = (rand.Float32() * depth) - (depth / 2.0)

		radius := 1.0 + (rand.Float32() * 4.0)

		sphere := raytracer.Sphere{
			Origin: pos,
			Radius: radius}
		
		sphere.Properties = getRandomMaterial()
		
		raytracer.Scene.AddObject(strconv.Itoa(int(i+1)), sphere)
	}
}

func getRandomMaterial() raytracer.Material {
	var material raytracer.Material
	color := getRandomColor()
	switch rand.Int() % 3 {
		case 0:
			material = raytracer.NewLambertian(color)
		case 1:
			material = raytracer.NewMetal(color, rand.Float32(), 1.0 + float32(int32(rand.Float32() * 300)), rand.Float32())
		case 2:
			material = raytracer.NewPhong(color, rand.Float32(), 1.0 + float32(int32(rand.Float32() * 70)))
	}
	
	return material
}

func getRandomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255}
}

func spawnSpheresInRing() {
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
}

func generatePhongShadingTest(numSpheresHigh, numSpheresWide uint) {
	reflectivity := float32(0.0)
	reflectivityInc := 1.0 / float32(numSpheresHigh)
	shininess := 1
	shininessInc := int(300.0 / float32(numSpheresWide))
	radius := float32(3.0)
	color := color.RGBA {
		R: 255,
		G: 255,
		B: 255,
		A: 255}
	sphereCount := 1
	for i := uint(0); i < numSpheresHigh; i++ {
		for j := uint(0); j < numSpheresWide; j++ {
			pos := raytracer.NewVector3(
				float32((radius * 2.5) * float32(j)),
				float32((radius * 2.5) * float32(i)),
				float32(-5.0))
			
			phong := raytracer.NewPhong(color, reflectivity, float32(shininess))
			
			sphere := raytracer.NewSphere(
				pos,
				radius,
				phong)
				
			raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
			sphereCount++
			shininess += shininessInc
		}
		
		shininess = 1
		reflectivity += reflectivityInc 
	}
	color.R = 0
	color.G = 0
	light := raytracer.NewPointLight(color, raytracer.NewVector3(52.5, 30, 100), 80)
	color.B = 0
	color.G = 255
	light2 := raytracer.NewPointLight(color, raytracer.NewVector3(70, 20, 100), 80)
	raytracer.Scene.AddLight("light", light)
	raytracer.Scene.AddLight("light2", light2)
	
	middleSphere := (int(numSpheresHigh / 2) * int(numSpheresWide)) + int(numSpheresWide / 2)
	fmt.Printf("Look up sphere %v", middleSphere)
}

func generateHeart() {
	color := color.RGBA {
		R: 255,
		A: 255 }
	sphereCount := 1
	material := raytracer.NewPhong(color, 0.3, 25)
	radius := float32(1.0)
	X := float32(0)
	Y := float32(0)
	Z := float32(-5.0)
	
	sphere := raytracer.NewSphere(
			raytracer.NewVector3(X, Y, Z),
			radius,
			material)
			
	raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
	sphereCount++
	
	for row := 0; row < 5; row++ {
		sphere1 := raytracer.NewSphere(
			raytracer.NewVector3(X, Y, Z),
			radius,
			material)
		sphere2 := raytracer.NewSphere(
			raytracer.NewVector3(-X, Y, Z),
			radius,
			material)
		
		X += radius * 2.0
		Y += radius * 2.0
		
		raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere1)
		sphereCount++
		
		raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere2)
		sphereCount++
	}
	
	sphere = raytracer.NewSphere(
			raytracer.NewVector3(0, 8, Z),
			radius,
			material)
	
	raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
	sphereCount++
	
	sphere = raytracer.NewSphere(
			raytracer.NewVector3(-6, 10, Z),
			radius,
			material)
	
	raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
	sphereCount++
	
	sphere = raytracer.NewSphere(
			raytracer.NewVector3(6, 10, Z),
			radius,
			material)
	
	raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
	sphereCount++
	
	sphere = raytracer.NewSphere(
			raytracer.NewVector3(4, 12, Z),
			radius,
			material)
	
	raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
	sphereCount++
	
	sphere = raytracer.NewSphere(
			raytracer.NewVector3(-4, 12, Z),
			radius,
			material)
	
	raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
	sphereCount++
	
	sphere = raytracer.NewSphere(
			raytracer.NewVector3(-2, 10, Z),
			radius,
			material)
	
	raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
	sphereCount++
	
	sphere = raytracer.NewSphere(
			raytracer.NewVector3(2, 10, Z),
			radius,
			material)
	
	raytracer.Scene.AddObject(strconv.Itoa(int(sphereCount)), sphere)
	sphereCount++
}

func generateCubes(height, width, depth float32, numBoxes uint) {
	for i := uint(0); i < numBoxes; i++ {
		var pos raytracer.Vector3
		pos.X = (rand.Float32() * width) - (width / 2.0)
		pos.Y = (rand.Float32() * height) - (height / 2.0)
		pos.Z = (rand.Float32() * depth) - (depth / 2.0)

		scaleX := rand.Float32() * 10.0
		// scaleY := rand.Float32() * 10.0
		// sacleZ := rand.Float32() * 10.0

		mesh := raytracer.NewExportableTriangleMesh("ObjectFiles\\bunny.obj", raytracer.NewVector3(scaleX, scaleX, scaleX), pos, nil)
		
		mesh.Properties = getRandomMaterial()
		
		raytracer.Scene.AddObject(strconv.Itoa(int(i+1)), mesh)
	}
}

func main() {
	var lightFilename string
	var sceneFilename string
	var numberSpheres uint 
	var numberLights uint
	var generatePhongTest bool
	var heartTest bool
	var sphereTest bool
	
	flag.StringVar(&lightFilename, "light", "GeneratedLights.json", "Filename to output generated lights JSON to")
	flag.StringVar(&sceneFilename, "scene", "GeneratedScene.json", "Filename to output the generated scene JSON to")
	flag.UintVar(&numberSpheres, "numSpheres", 50, "The number of spheres to spawn")
	flag.UintVar(&numberLights, "numLights", 5, "Number of lights to spawn")
	flag.BoolVar(&generatePhongTest, "phongtest", false, "Pass true to this to generate a phong material test grid")
	flag.BoolVar(&heartTest, "heart", false, "Pass true to this to generate a heart")
	flag.BoolVar(&sphereTest, "spherebox", false, "Pass true to this to generate spheres in a box")
	flag.Parse()	
	
	seedVal := time.Now().UTC().UnixNano()
	rand.Seed(seedVal)
	fmt.Printf("Generating with seed val: %v\n", seedVal)
	
	if (generatePhongTest) {
		fmt.Println("Generating phong test grid")
		generatePhongShadingTest(20, 30)
	} else if (heartTest) {
		generateHeart()
		spawnLightsInXZPlane(40.0, 3)	
	} else if (sphereTest) {
		spawnSpheresInCube(100.0, 100.0, 100.0, numberSpheres)
		spawnLightsInXZPlane(200.0, numberLights)	
	} else {
		generateCubes(100.0, 100.0, 100.0, numberSpheres)
		spawnLightsInXZPlane(200.0, numberLights)
	}
	
	raytracer.ExportScene(sceneFilename, lightFilename)
}