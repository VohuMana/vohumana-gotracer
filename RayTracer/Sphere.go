package raytracer

import
(
    "image/color"
    "math"
    "math/rand"
)

// Sphere is basic geometry used by the ray tracer
type Sphere struct {
    Origin Vector3
    Radius float32 
    Properties Material
}

func createUnitSphereVector() Vector3 {
    return Vector3 {
        rand.Float32(),
        rand.Float32(),
        rand.Float32() }.Subtract(Vector3 {
            1.0,
            1.0,
            1.0 }).Multiply(Vector3 {
                2.0,
                2.0,
                2.0 })
}

func  randomVectorInUnitSphere() Vector3 {
    p := createUnitSphereVector()
    
    for p.Dot(p) >= 1.0 {
        p = createUnitSphereVector()
    }
    
    return p
}

// TestIntersection will test for an intersection between the sphere and ray
func (s Sphere) TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    var record IntersectionRecord
    
    m := r.Origin.Subtract(s.Origin)
    b := m.Dot(r.Direction)
    c := m.Dot(m) - (s.Radius * s.Radius)
    
    if (c > 0.0 && b > 0.0) {
        return false, record
    }
    
    descriminant := (b * b) - c
    
    if (descriminant < 0.0) {
        return false, record
    }
    
    record.T = float32(float64(-b) - math.Sqrt(float64(descriminant)))
    
    if (record.T < tMin || record.T > tMax) {
        return false, record
    }
    
    record.Point = r.PointOnRay(record.T)
    record.Normal = record.Point.Subtract(s.Origin).UnitVector()
    record.Object = s
    
    return true, record
}

// GetColor gets the color at a collision point
func (s Sphere) GetColor(i IntersectionRecord, bounces uint32) color.RGBA {    
    if (bounces > MaxBounces) {
        return s.Properties.Color
    }
    
    target := i.Point.Add(i.Normal).Add(randomVectorInUnitSphere())
    bouncedRay := Ray {
        Origin: i.Point,
        Direction: target.Subtract(i.Point).UnitVector() }
    c := ShootRay(bouncedRay, Scene, bounces + 1)
    c.R = uint8(restrictValues(s.Properties.Reflectiveness * float32(c.R), 0.0, 255.0))
    c.G = uint8(restrictValues(s.Properties.Reflectiveness * float32(c.G), 0.0, 255.0))
    c.B = uint8(restrictValues(s.Properties.Reflectiveness * float32(c.B), 0.0, 255.0))
    
    c.R = uint8(restrictValues(float32(c.R + s.Properties.Color.R), 0.0, 255.0))
    c.G = uint8(restrictValues(float32(c.G + s.Properties.Color.G), 0.0, 255.0))
    c.B = uint8(restrictValues(float32(c.B + s.Properties.Color.B), 0.0, 255.0))
    return c
    
    // Render normals
    // c := i.Normal.Scale(0.5)
    // c = c.Add(Vector3{X:1.0, Y:1.0, Z:1.0})
    
    // c.X = float32(math.Min(1.0, float64(c.X)))
    // c.Y = float32(math.Min(1.0, float64(c.Y)))
    // c.Z = float32(math.Min(1.0, float64(c.Z)))
    // return c.AsColor()
}

func restrictValues(num, min, max float32) float32 {
    num = float32(math.Min(float64(num), float64(max)))
    num = float32(math.Max(float64(num), float64(min)))
    return num
}