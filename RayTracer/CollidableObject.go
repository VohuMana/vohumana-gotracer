package raytracer

import
(
    "image/color"
    "math"
)

// CollidableObject is an interface for objects that want to be able to collide with rays
type CollidableObject interface {
    TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord)
    GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA
}

// IntersectionRecord is an object that contains data about where a ray hit an object
type IntersectionRecord struct {
    T float32
    Point Vector3
    Normal Vector3
    Object CollidableObject
}

// calculateReflectionVector calulates a reflection vector with direction d and normal n. r = d - 2(n dot d)*n
func calculateReflectionVector(d, n Vector3) Vector3 {
    return d.Subtract(n.Scale(2.0 * d.Dot(n))).UnitVector()
}

// getNormalAsColor give a normal n, return the color value for that normal
func getNormalAsColor(n Vector3) color.RGBA {
    // Render normals
    c := n.Scale(0.5)
    c = c.Add(Vector3{X:1.0, Y:1.0, Z:1.0})
    
    c.X = float32(math.Min(1.0, float64(c.X)))
    c.Y = float32(math.Min(1.0, float64(c.Y)))
    c.Z = float32(math.Min(1.0, float64(c.Z)))
    return c.AsColor()
}