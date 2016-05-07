package raytracer

import
(
    "encoding/json"
    "image/color"
    "log"
    "math"
    "os"
)

// World contains information about the world
type World struct {
    Scene CollisionList
}

// Config contains data on how the raytracer will behave
type Config struct {
    // SkyColorTop is the color of the sky at the top of the picture
    SkyColorTop Vector3

    // SkyColorBottom is the color of the sky at the bottom of the picture
    SkyColorBottom Vector3

    // MaxBounces is the maximum number of bounces that can occur before the ray tracer stops reflecting rays
    MaxBounces uint32

    // MaxRaysPerBounce is the maximum number of rays that will be bounced from a single intersection
    MaxRaysPerBounce uint32

    // MaxAntialiasRays is the maximum number of rays that will be shot per pixel for anit aliasing
    MaxAntialiasRays uint32
}

// Settings contains the current config the ray tracer will use
var Settings Config 

// Scene is the global variable that should contain the scene
var Scene World

// AddObject adds a collidableobject to the scene
func (w *World) AddObject(name string, obj CollidableObject) {
    w.Scene.addObject(name, obj)
}

// TestCollision tests all the objects in the scene for collisions
func (w World) TestCollision(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    return w.Scene.testCollision(r, tMin, tMax)
}

// ShootRay shoots a ray and tests for intersection
func ShootRay(r Ray, w World, bounceDepth uint32) color.RGBA {
    collided, record := w.TestCollision(r, 0.0001, math.MaxFloat32)
    if (collided) {
        return record.Object.GetColor(r, record, bounceDepth)
    }
    
    t := 0.5 * (r.Direction.Y + 1.0)
    // Lerp from blue to white
    c := Settings.SkyColorBottom.Scale(1.0 - t).Add(Settings.SkyColorTop.Scale(t)); 
    return c.AsColor()
}

func checkError(err error) {
	if (err != nil) {
		log.Fatal(err)
    }
}

func ExportScene(filename string) {
    sceneString, err := json.Marshal(Scene.Scene.collisionList)
    checkError(err)
    
    jsonString := []byte{}
    jsonString = append(jsonString, sceneString...)
    
    sceneFile, err := os.Create(filename)
    checkError(err)
    defer sceneFile.Close()
    
    sceneFile.Write(jsonString)
}

func ExportConfig(filename string) {
    configString, err := json.Marshal(Settings)
    checkError(err)
    
    configFile, err := os.Create(filename)
    checkError(err)
    defer configFile.Close()
    
    configFile.Write(configString)
}