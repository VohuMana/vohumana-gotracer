package raytracer

import
(
    "image/color"
)

// Material is a struct that describes the properties each object should have
type Material struct {
    Color color.RGBA
    Reflectiveness float32
    IsDiffuse bool
}