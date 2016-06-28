package raytracer

import 
(
    "math"
)

// BoundingBox is a box that encapsulates all points within it
type BoundingBox struct {
    Min, Max Vector3
}

func minVectorValues(a, b Vector3) Vector3 {
    return NewVector3(
        float32(math.Min(float64(a.X), float64(b.X))),
        float32(math.Min(float64(a.Y), float64(b.Y))),
        float32(math.Min(float64(a.Z), float64(b.Z))))
}

func maxVectorValues(a, b Vector3) Vector3 {
    return NewVector3(
        float32(math.Max(float64(a.X), float64(b.X))),
        float32(math.Max(float64(a.Y), float64(b.Y))),
        float32(math.Max(float64(a.Z), float64(b.Z))))
}

// GenerateBoundingBoxFromTris will create a bounding box from triangles
func GenerateBoundingBoxFromTris(tris []Triangle) BoundingBox {
    min := NewVector3(math.MaxFloat32, math.MaxFloat32, math.MaxFloat32)
    max := NewVector3(-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32)

    for _, tri := range tris {
        min = minVectorValues(min, tri.Position1)
        max = maxVectorValues(max, tri.Position1)

        min = minVectorValues(min, tri.Position2)
        max = maxVectorValues(max, tri.Position2)

        min = minVectorValues(min, tri.Position3)
        max = maxVectorValues(max, tri.Position3)
    }

    return BoundingBox {
        Min: min, 
        Max: max }
}

// IsRayColliding tests for ray AABB collision but does not care about the point of the collision
func (b BoundingBox) IsRayColliding(r Ray) bool {
    return false
}