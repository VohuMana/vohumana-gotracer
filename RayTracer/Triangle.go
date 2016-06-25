package raytracer

import
(
    "image/color"
)

// Triangle is a object used to generate triangle meshes
type Triangle struct {
    Position1, Position2, Position3, Normal Vector3
    Properties Material
    Midpoint Vector3
}

// NewTriangle creates a new triangle with given properties
func NewTriangle(pos1, pos2, pos3 Vector3, mat Material) Triangle {
    return NewTriangleWithNormal(pos1, pos2, pos3, pos2.Subtract(pos1).Cross(pos3.Subtract(pos1)).UnitVector(), mat)
}

// NewTriangleWithNormal will create a triangle with a specified normal
func NewTriangleWithNormal(pos1, pos2, pos3, norm Vector3, mat Material) Triangle {
    return Triangle {
        Position1: pos1,
        Position2: pos2,
        Position3: pos3,
        Normal: norm,
        Midpoint: pos1.Add(pos2.Add(pos3)).Scale(1.0 / 3.0),
        Properties: mat }
}

// TestIntersection will test if a ray is colliding with the triangle
func (tri Triangle) TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    edge1 := tri.Position2.Subtract(tri.Position1)
    edge2 := tri.Position3.Subtract(tri.Position1)

    pvec := r.Direction.Cross(edge2)
    determinant := pvec.Dot(edge1)
    
    if (determinant < 0.0001) {
        return false, IntersectionRecord{}
    }

    tvec := r.Origin.Subtract(tri.Position1)
    u := tvec.Dot(pvec)

    if (u < 0.0 || u > determinant) {
        return false, IntersectionRecord{}
    }

    qvec := tvec.Cross(edge1)
    v := r.Direction.Dot(qvec)

    if (v < 0.0 || (u + v) > determinant) {
        return false, IntersectionRecord{}
    }

    inverseDeterminant := 1.0 / determinant
    t := edge2.Dot(qvec) * inverseDeterminant

    if (t > tMin && t < tMax) {
        return true, IntersectionRecord {
        T: t,
        Normal: tri.Normal,
        Point: r.PointOnRay(t),
        Object: tri }
    }

    return false, IntersectionRecord{}
}

// GetColor will get the color at the collision point
func (tri Triangle) GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA {
    return tri.Properties.GetColor(r, i, Scene, bounces)
}

// // ApplyMatrix3 will apply a matrix to the triangle
// func (tri Triangle) ApplyMatrix3(mtx Matrix3) Triangle {    
//     return NewTriangle(
//         mtx.MultiplyVector3(tri.Position1),
//         mtx.MultiplyVector3(tri.Position2),
//         mtx.MultiplyVector3(tri.Position3),
//         tri.Properties)
// }

// // Translate will translate the triangle
// func (tri Triangle) Translate(translation Vector3) Triangle {
//     return NewTriangle(
//         tri.Position1.Add(translation),
//         tri.Position2.Add(translation),
//         tri.Position3.Add(translation),
//         tri.Properties)
// }