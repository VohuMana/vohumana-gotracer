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

func restrictValues(num, min, max float32) float32 {
    num = float32(math.Min(float64(num), float64(max)))
    num = float32(math.Max(float64(num), float64(min)))
    return num
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

func calculateReflectionRay(r Ray, i IntersectionRecord) Ray {
    return Ray {
        Origin: i.Point,
        Direction: calculateReflectionVector(r.Direction, i.Normal) }
}

func calculateDiffuseRay(i IntersectionRecord) Ray {
        target := i.Point.Add(i.Normal).Add(randomVectorInUnitSphere())
        return Ray {
            Origin: i.Point,
            Direction: target.Subtract(i.Point).UnitVector() }
}

// GetColor gets the color at a collision point
func (s Sphere) GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA {    
    // If the ray has bounced more times than the provided amout return this objects' color
    if (bounces > MaxBounces) {
        return s.Properties.Color
    }
    
    var bouncedRay Ray
    var c color.RGBA
    
    if (s.Properties.IsDiffuse) {
        // Calculate multiple scattered rays
        bouncedRay = calculateDiffuseRay(i)
        var r, g, b float32
        
        // Bounce multiple diffuse rays
        for i := uint32(0); i < MaxDiffuseRays; i++ {
            color := ShootRay(bouncedRay, Scene, bounces + 1)
            r += float32(color.R)
            g += float32(color.G)
            b += float32(color.B)
        }
        
        // Average the color
        r /= float32(MaxDiffuseRays)
        g /= float32(MaxDiffuseRays)
        b /= float32(MaxDiffuseRays)
        
        // Set the averaged color
        c.R = uint8(r)
        c.G = uint8(g)
        c.B = uint8(b)
        c.A = uint8(255)
        
    } else {
        // Calculate where this ray would bounce to
        bouncedRay = calculateReflectionRay(r, i)
        
        // Shoot the reflected ray out into the scene and see where it collides with
        c = ShootRay(bouncedRay, Scene, bounces + 1)
    }
    
    // Calculate the returned color value after applying how difuse the material is
    c.R = uint8(restrictValues(s.Properties.Reflectiveness * float32(c.R), 0.0, 255.0))
    c.G = uint8(restrictValues(s.Properties.Reflectiveness * float32(c.G), 0.0, 255.0))
    c.B = uint8(restrictValues(s.Properties.Reflectiveness * float32(c.B), 0.0, 255.0))
    
    // Add this objects' color value to the result    
    c.R = uint8(restrictValues(float32(c.R) + (float32(s.Properties.Color.R) * (1.0 - s.Properties.Reflectiveness)), 0.0, 255.0))
    c.G = uint8(restrictValues(float32(c.G) + (float32(s.Properties.Color.G) * (1.0 - s.Properties.Reflectiveness)), 0.0, 255.0))
    c.B = uint8(restrictValues(float32(c.B) + (float32(s.Properties.Color.B) * (1.0 - s.Properties.Reflectiveness)), 0.0, 255.0))
    return c
}