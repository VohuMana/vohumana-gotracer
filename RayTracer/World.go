package raytracer

import
(
    "encoding/json"
    "image/color"
    "io"
    "log"
    "math"
    "os"
)

// World contains information about the world
type World struct {
    Scene CollisionList
    Lights LightList
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
    
    // MaxLightRays is the maximum number of rays that will be shot per light for smooth shadowing
    MaxLightRays uint32
    
    // WidthInPixels is the horizontal resolution of the resulting image
    WidthInPixels int
    
    // HeightInPixels is the vertical resolution of the resulting image    
    HeightInPixels int
}

// Settings contains the current config the ray tracer will use
var Settings Config 

// Scene is the global variable that should contain the scene
var Scene World

// AddObject adds a collidableobject to the scene
func (w *World) AddObject(name string, obj CollidableObject) {
    w.Scene.addObject(name, obj)
}

// AddLight adds a light to the world
func (w *World) AddLight(name string, light Light) {
    w.Lights.addObject(name, light)
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

// ExportScene will export the current global scene
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

// ImportScene will import the given scene file
func ImportScene(filename string) {
    sceneFile, err := os.Open(filename)
    checkError(err)
    defer sceneFile.Close()
    
    info, err := sceneFile.Stat()
    checkError(err)
    
    contents := make([]byte, info.Size())
    
    _, err = sceneFile.Read(contents)
    if (err != nil && io.EOF != err) {
        checkError(err)   
    }
    
    var sceneObjects map[string]interface{}
    err = json.Unmarshal(contents, &sceneObjects)
    checkError(err)
    
    for name, object := range sceneObjects {
        switch object.(type) {
            case map[string]interface{}:
                obj, ok := object.(map[string]interface{})
                if (true == ok) {
                    s, isSphere := deserializeSphere(obj)
                    if (true == isSphere) {
                        Scene.AddObject(name, s)
                    }    
                }
            default:
                continue
        }
    }    
}

// ExportConfig will export the current global config
func ExportConfig(filename string) {
    configString, err := json.Marshal(Settings)
    checkError(err)
    
    configFile, err := os.Create(filename)
    checkError(err)
    defer configFile.Close()
    
    configFile.Write(configString)
}

// ImportConfig will import a config file to the global config
func ImportConfig(filename string) {
    configFile, err := os.Open(filename)
    checkError(err)
    defer configFile.Close()
    
    info, err := configFile.Stat()
    checkError(err)
    
    contents := make([]byte, info.Size())
    
    _, err = configFile.Read(contents)
    if (err != nil && io.EOF != err) {
        checkError(err)   
    }
    
    err = json.Unmarshal(contents, &Settings)
    checkError(err)
}