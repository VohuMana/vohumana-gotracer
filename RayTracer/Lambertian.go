package raytracer

import (
	"image/color"
)

// Lambertian is a type of material that scatters rays randomly, used for diffuse objectss
type Lambertian struct {
	BaseProperties Phong
}

// NewLambertian will return a new lambertian material
func NewLambertian(color color.RGBA) Lambertian {
	return Lambertian{
		BaseProperties: NewPhong(color, 0.0, 1)}
}

// GetColor for lambertian materials
func (l Lambertian) GetColor(r Ray, i IntersectionRecord, w World, bounceDepth uint32) color.RGBA {
	if bounceDepth > Settings.MaxBounces {
		return l.BaseProperties.DiffuseColor.AsColor()
	}

	diffuseRay := calculateDiffuseRay(i)
	return calculatePhongLighting(l.BaseProperties, i, diffuseRay, w, bounceDepth)
}
