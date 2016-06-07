package raytracer

import
(
    "image/color"
)

// Triangle is a object used to generate triangle meshes
type Triangle struct {
    Position1, Position2, Position3, Normal Vector3
    Properties Material
}

// NewTriangle creates a new triangle with given properties
func NewTriangle(pos1, pos2, pos3 Vector3, mat Material) Triangle {
    return Triangle {
        Position1: pos1,
        Position2: pos2,
        Position3: pos3,
        Normal: pos2.Subtract(pos1).Cross(pos3.Subtract(pos1)).UnitVector(),
        Properties: mat }
}

// TestIntersection will test if a ray is colliding with the triangle
func (tri Triangle) TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    a := tri.Position2.Subtract(tri.Position1)
    b := tri.Position3.Subtract(tri.Position1)
    
    t := r.Origin.Subtract(tri.Position1).Scale(1.0).Dot(tri.Normal) / r.Direction.Dot(tri.Normal)
    
    if (t <= tMin || t >= tMax) {
        return false, IntersectionRecord{}
    }
    
    aDotB := a.Dot(b)
    coeff := float32(1.0) / (float32(a.SquareLength()) * float32(b.SquareLength()) - (aDotB * aDotB))
    intersectionPoint := r.PointOnRay(t)
    inter := intersectionPoint.Subtract(tri.Position1)
    interDotA := inter.Dot(a)
    interDotB := inter.Dot(b)
    alpha := coeff * (interDotA * float32(b.SquareLength()) + -aDotB * interDotB)
    beta := coeff * (interDotA * -aDotB + float32(a.SquareLength()) * interDotB)
    
    if (alpha < 0.0 || beta < 0.0 || (alpha + beta) > 1.0) {
        return false, IntersectionRecord{}
    }
    
    return true, IntersectionRecord {
        T: t,
        Normal: tri.Normal,
        Point: intersectionPoint,
        Object: tri }
}

// GetColor will get the color at the collision point
func (tri Triangle) GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA {
    return tri.Properties.GetColor(r, i, Scene, bounces)
}