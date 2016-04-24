package raytracer

import
(
    "image/color"
    "math"
)

// World contains information about the world
type World struct {
    Scene CollisionList
}

// SkyColorTop is the color of the sky at the top of the picture
var SkyColorTop Vector3

// SkyColorBottom is the color of the sky at the bottom of the picture
var SkyColorBottom Vector3

// Scene is the global variable that should contain the scene
var Scene World

// MaxBounces is the maximum number of bounces that can occur before the ray tracer stops reflecting rays
var MaxBounces uint32

// MaxRaysPerBounce is the maximum number of rays that will be bounced from a single intersection
var MaxRaysPerBounce uint32

// MaxAntialiasRays is the maximum number of rays that will be shot per pixel for anit aliasing
var MaxAntialiasRays uint32

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
    c := SkyColorBottom.Scale(1.0 - t).Add(SkyColorTop.Scale(t)); 
    return c.AsColor()
}