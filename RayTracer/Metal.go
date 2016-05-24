package raytracer

import
(
    "image/color"
)

// Metal is a type of material that reflects rays
type Metal struct {
    Fuzziness float32
    BaseProperties Phong
}

// NewMetal will create a new metal object with the given properties
func NewMetal(color color.RGBA, reflectivity, shininess, fuzziness float32) Metal {
    return Metal {
        Fuzziness: fuzziness,
        BaseProperties: NewPhong(color, reflectivity, shininess) }
}

// GetColor calculates the color for a given intersection with this material
func (m Metal) GetColor(r Ray, i IntersectionRecord, w World, bounceDepth uint32) color.RGBA {
    if (bounceDepth > Settings.MaxBounces) {
        return m.BaseProperties.DiffuseColor.AsColor()
    }
    
    reflectedRay := calculateReflectionRay(r, i, m.Fuzziness)
    return calculatePhongLighting(m.BaseProperties, i, reflectedRay, w, bounceDepth)
}