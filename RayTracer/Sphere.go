package raytracer

// Sphere is basic geometry used by the ray tracer
type Sphere struct {
    Origin Vector3
    Radius float32 
}

// TestIntersection will test for an intersection between the sphere and ray
func (s Sphere) TestIntersection(r Ray) bool {
    return false
}