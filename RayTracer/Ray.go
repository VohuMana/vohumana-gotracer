package raytracer

// Ray is a mathematical ray having a starting point and direction vector
type Ray struct {
    Origin, Direction Vector3
}

// PointOnRay will get a point on the ray at time t Origin + (t * Direction)
func (r Ray) PointOnRay(t float32) Vector3 {
    return r.Origin.Add(r.Direction.Scale(t))
}