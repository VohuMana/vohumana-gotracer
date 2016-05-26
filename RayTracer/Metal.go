package raytracer

import (
	"image/color"
)

// Metal is a type of material that reflects rays
type Metal struct {
	Fuzziness      float32
	BaseProperties Phong
}

// NewMetal will create a new metal object with the given properties
func NewMetal(color color.RGBA, reflectivity, shininess, fuzziness float32) Metal {
	return Metal{
		Fuzziness:      fuzziness,
		BaseProperties: NewPhong(color, reflectivity, shininess)}
}

// GetColor calculates the color for a given intersection with this material
func (m Metal) GetColor(r Ray, i IntersectionRecord, w World, bounceDepth uint32) color.RGBA {
	if bounceDepth > Settings.MaxBounces {
		return m.BaseProperties.DiffuseColor.AsColor()
	}

	var c color.RGBA
	var red, green, blue float32

	// Bounce multiple diffuse rays
	for rays := uint32(0); rays < Settings.MaxRaysPerBounce; rays++ {
		reflectedRay := calculateReflectionRay(r, i, m.Fuzziness)
		return calculatePhongLighting(m.BaseProperties, i, reflectedRay, w, bounceDepth)
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

	return c
}
