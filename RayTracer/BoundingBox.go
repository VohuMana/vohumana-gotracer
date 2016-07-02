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
    dirFrac := NewVector3(1.0, 1.0, 1.0)
    dirFrac = dirFrac.Divide(r.Direction)

    t1 := float64((b.Min.X - r.Origin.X) * dirFrac.X)
    t2 := float64((b.Max.X - r.Origin.X) * dirFrac.X)
    t3 := float64((b.Min.Y - r.Origin.Y) * dirFrac.Y)
    t4 := float64((b.Max.Y - r.Origin.Y) * dirFrac.Y)
    t5 := float64((b.Min.Z - r.Origin.Z) * dirFrac.Z)
    t6 := float64((b.Max.Z - r.Origin.Z) * dirFrac.Z)

    tMin := math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
    tMax := math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

    if (tMax < 0) {
        return false
    }

    if (tMin > tMax) {
        return false
    }

    return true
}