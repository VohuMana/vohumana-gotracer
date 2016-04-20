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

func calculateReflectionRay(r Ray, i IntersectionRecord, fuzziness float32) Ray {
    return Ray {
        Origin: i.Point,
        Direction: calculateReflectionVector(r.Direction, i.Normal).Add(randomVectorInUnitSphere().Scale(fuzziness)) }
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
        return color.RGBA {
            R: 0,
            G: 0,
            B: 0,
            A: 255 }
    }
    
    var bouncedRay Ray
    var c color.RGBA
    var red, green, blue float32
    
    // Bounce multiple diffuse rays
    for rays := uint32(0); rays < MaxRaysPerBounce; rays++ {
        if (s.Properties.IsDiffuse) {
            // Calculate multiple scattered rays
            bouncedRay = calculateDiffuseRay(i)
        } else {
            // Calculate where this ray would bounce to
            bouncedRay = calculateReflectionRay(r, i, s.Properties.Fuzziness)
        }
        
        color := ShootRay(bouncedRay, Scene, bounces + 1)
        red += float32(color.R)
        green += float32(color.G)
        blue += float32(color.B)
    }
    
    // Average the color
    red /= float32(MaxRaysPerBounce)
    green /= float32(MaxRaysPerBounce)
    blue /= float32(MaxRaysPerBounce)
    
    // Set the averaged color
    c.R = uint8(red)
    c.G = uint8(green)
    c.B = uint8(blue)
    c.A = uint8(255)        
    
    // Multiply this objects color with the incoming color
    c = AsVector3(c).Multiply(s.Properties.Attenuation).AsColor()
    
    return c
}