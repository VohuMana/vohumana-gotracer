package raytracer

import
(
    "image/color"
)

// Emissive is a type of material that glows
type Emissive struct {
    Emission Vector3
}

// NewEmissive returns a new emissive material with given value
func NewEmissive(emissionColor Vector3) Emissive {
    return Emissive {
        Emission: emissionColor }
}

// GetColor calculates the color for a given intersection with this material
func (e Emissive) GetColor(r Ray, i IntersectionRecord, w World, bounceDepth uint32) color.RGBA {
    // Emissive materials are not affected by lights nor do they reflect anything
    return e.Emission.AsColor()
}