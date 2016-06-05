package raytracer

import
(
    "image/color"
)

// InfinitePlane is a 3D plane that extends in all directions infinitly
type InfinitePlane struct {
    Position, Normal Vector3
    Properties Material
}

// NewInfinitePlane creates a new infinite plane
func NewInfinitePlane(pos, normal Vector3, mat Material) InfinitePlane {
    return InfinitePlane {
        Position: pos,
        Normal: normal,
        Properties: mat }
}

// TestIntersection will test if a ray is colliding with the infinite plane
func (p InfinitePlane) TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    distanceVec := p.Position.Subtract(r.Origin)
    nearestDistance := p.Normal.Dot(distanceVec)
    speed := r.Direction.Dot(p.Normal)
    t := nearestDistance / speed
    
    if (t <= 0.0 || t <= tMin || t >= tMax) {
        return false, IntersectionRecord{}
    }
    
    return true, IntersectionRecord {
        T: t,
        Normal: p.Normal,
        Point: r.PointOnRay(t),
        Object: p }
}

// GetColor will get the color at the collision point
func (p InfinitePlane) GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA {
    return p.Properties.GetColor(r, i, Scene, bounces)
}