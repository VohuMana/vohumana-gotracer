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

// Scene is the global variable that should contain the scene
var Scene World

var MaxBounces uint32

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
    collided, record := w.TestCollision(r, 0.0, math.MaxFloat32)
    if (collided) {
        return record.Object.GetColor(record, bounceDepth)
    }
    
    t := 0.5 * (r.Direction.Y + 1.0)
    // Lerp from blue to white
    c := Vector3{ X: 0.5, Y: 0.7, Z: 1.0 }.Scale(1.0 - t).Add(Vector3{X:1.0, Y: 1.0, Z: 1.0}.Scale(t)); 
    return c.AsColor()
}