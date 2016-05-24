package raytracer

import
(
    "image/color"
)

// Phong is a type of material that can be used in lighting
type Phong struct {
    DiffuseColor Vector3
    Shininess float32
    Reflectivity float32
}

// NewPhong creates a new phong material object
func NewPhong(color color.RGBA, reflectivity, shininess float32) Phong {
    return Phong {
        DiffuseColor: AsVector3(color),
        Shininess: shininess,
        Reflectivity: reflectivity }
}

// GetColor calculates the color for a given intersection with this material
func (p Phong) GetColor(r Ray, i IntersectionRecord, w World, bounceDepth uint32) color.RGBA {
    if (bounceDepth > Settings.MaxBounces) {
        return p.DiffuseColor.AsColor()
    }
    
    reflectedRay := calculateReflectionRay(r, i, 0.0)
    return calculatePhongLighting(p, i, reflectedRay, w, bounceDepth)
}