package raytracer

import
(
    "image/color"
    "math"
)

// Sphere is basic geometry used by the ray tracer
type Sphere struct {
    Origin Vector3
    Radius float32 
    Properties Material
}

// TestIntersection will test for an intersection between the sphere and ray
func (s Sphere) TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    var record IntersectionRecord
    
    // Make a vector from the sphere origin to the ray origin
    m := r.Origin.Subtract(s.Origin)
    
    // Dot the direction of the ray and the direction of m.  They must face opposite ways for their to be collision, ie they must have a 0 or negative dot product.
    b := m.Dot(r.Direction)
    
    // Subtract the length squared of m and R^2, they should be 0 or less for collision
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
func (s Sphere) GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA {    
    // If the ray has bounced more times than the provided amout return this objects' color
    if (bounces > Settings.MaxBounces) {
        return color.RGBA {
            R: 255,
            G: 255,
            B: 255,
            A: 255 }
    }
    
    var bouncedRay Ray
    var c color.RGBA
    
    if (false == s.Properties.IsEmissive()) {
        var red, green, blue float32
    
        // Bounce multiple diffuse rays
        for rays := uint32(0); rays < Settings.MaxRaysPerBounce; rays++ {
            bouncedRay = s.Properties.Scatter(r, i)
            
            color := ShootRay(bouncedRay, Scene, bounces + 1)
            red += float32(color.R)
            green += float32(color.G)
            blue += float32(color.B)
        }
        
        // Average the color
        red /= float32(Settings.MaxRaysPerBounce)
        green /= float32(Settings.MaxRaysPerBounce)
        blue /= float32(Settings.MaxRaysPerBounce)
        
        // Set the averaged color
        c.R = uint8(red)
        c.G = uint8(green)
        c.B = uint8(blue)
        c.A = uint8(255)        
        
        // Multiply this objects color with the incoming color
        c = AsVector3(c).Multiply(s.Properties.GetAttenuation()).AsColor()    
    } else {
        c = s.Properties.GetEmission().AsColor()
    }
    
    return c
}

func deserializeSphere(object map[string]interface{}) (Sphere, bool) {
    var sphere Sphere
    validSphere := true
    
    for name, object := range object {
        switch name {
            case "Origin":
                origin, ok := object.(map[string]interface{})
                if (true == ok) {
                    sphere.Origin, ok = deserializeVector3(origin)
                    if (false == ok) {
                        validSphere = false
                        break;
                    }
                }
                
            case "Radius":
                radius, ok := object.(float64)
                if (true == ok) {
                    sphere.Radius = float32(radius)
                } else {
                    validSphere = false
                    break;
                }
                
            case "Properties":
                prop, ok := object.(map[string]interface{})
                if (true == ok) {
                    mat, ok := deserializeMaterial(prop)
                    if (true == ok) {
                        sphere.Properties = mat
                    } else {
                        validSphere = false
                        break;
                    }    
                }
            
            default:
                validSphere = false
                break
        }
    }
    
    return sphere, validSphere
}