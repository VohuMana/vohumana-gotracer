package raytracer

import
(
    "image/color"
)

// Dielectric is a type of material that refracts rays
type Dielectric struct {
    RefractiveIndex float32
    Reflectivity float32
    Shininess float32
}

// NewDielectric creates a new dielectric material with the given parameters
func NewDielectric(reflectivity, shininess, refractiveIndex float32) Dielectric {
    return Dielectric {
        RefractiveIndex: refractiveIndex,
        Reflectivity: reflectivity,
        Shininess: shininess }
}

// GetColor calculates the color for a given intersection with this material
func (d Dielectric) GetColor(r Ray, i IntersectionRecord, w World, bounceDepth uint32) color.RGBA {
    if (bounceDepth > Settings.MaxBounces) {
        return color.RGBA {
            R: 255,
            G: 255,
            B: 255,
            A: 255 }
    }
    
    // Base color of the dielectric will be the refracted color
    var c color.RGBA
    var red, green, blue float32
    
    // Bounce multiple diffuse rays
    for rays := uint32(0); rays < Settings.MaxRaysPerBounce; rays++ {
        refractedRay := calculateRefractedRay(r, i, d.RefractiveIndex)
        
        color := ShootRay(refractedRay, Scene, bounceDepth + 1)
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
    
    phong := NewPhong(c, d.Reflectivity, d.Shininess)
    reflectedRay := calculateReflectionRay(r, i, 0.0)
    
    return calculatePhongLighting(phong, i, reflectedRay, w, bounceDepth)
}