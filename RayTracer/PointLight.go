package raytracer

import (
	"image/color"
)

// PointLight is a light that emits light in all directions
type PointLight struct {
	Color, Position Vector3
	Power           float32
}

// NewPointLight returns a new point light with the given values
func NewPointLight(color color.RGBA, position Vector3, power float32) PointLight {
	return PointLight{
		Color:    AsVector3(color),
		Position: position,
		Power:    power}
}

// GetColor returns the color of the light
func (p PointLight) GetColor() Vector3 {
	return p.Color
}

// GetPosition returns the current position of the light
func (p PointLight) GetPosition() Vector3 {
	return p.Position
}

// GetDirection point light returns false because light is emitted in all directions evenly
func (p PointLight) GetDirection() (Vector3, bool) {
	return NewVector3(0, 0, 0), false
}

// GetPower will return the lights power
func (p PointLight) GetPower() float32 {
	return p.Power
}

// CalculateColor will calculate the color of the intersection
func (p PointLight) CalculateColor(r Ray, material Phong, i IntersectionRecord) Vector3 {
	return NewVector3(0, 0, 0)
}
