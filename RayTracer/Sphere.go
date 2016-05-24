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

// NewSphere returns a new sphere with the given values
func NewSphere(origin Vector3, radius float32, properties Material) Sphere {
    return Sphere {
        Origin: origin,
        Radius: radius,
        Properties: properties }
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
    return s.Properties.GetColor(r, i, Scene, bounces)
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