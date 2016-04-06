package raytracer

import
(
    "image/color"
    "math"
)

// Vector3 is a basic 3D vector type with members X, Y, and Z
type Vector3 struct {
    X, Y, Z float32
}

// Add adds vectors a and b like a + b
func (v Vector3) Add(a Vector3) Vector3 {
    return Vector3 {
        X: v.X + a.X,
        Y: v.Y + a.Y,
        Z: v.Z + a.Z }
}

// Subtract subtracts vectors a and b like a - b
func (v Vector3) Subtract(a Vector3) Vector3 {
    return Vector3 {
        X: v.X - a.X,
        Y: v.Y - a.Y,
        Z: v.Z - a.Z }
}

// Multiply multiplies vectors a and b like a * b
func (v Vector3) Multiply(a Vector3) Vector3  {
    return Vector3 {
        X: v.X * a.X,
        Y: v.Y * a.Y,
        Z: v.Z * a.Z }
}

// Scale scales a vector by factor s
func (v Vector3) Scale(s float32) Vector3 {
    return Vector3 {
        X: v.X * s,
        Y: v.Y * s,
        Z: v.Z * s }
}

// Divide divides vectors a and b like a / b
func (v Vector3) Divide(a Vector3) Vector3 {
    return Vector3 {
        X: v.X / a.X,
        Y: v.Y / a.Y,
        Z: v.Z / a.Z }
}

// Length will get the length of the vector as defined bY sqrt(X^2 + Y^2 + Z^2)
func (v Vector3) Length() float64 {
    return math.Sqrt(float64(v.X * v.X) + float64(v.Y * v.Y) + float64(v.Z * v.Z))
}

// SquareLength will return the length of the vector squared.  X^2 + Y^2 + Z^2
func (v Vector3) SquareLength() float64 {
    return float64(v.X * v.X) + float64(v.Y * v.Y) + float64(v.Z * v.Z)
}

// Dot will return the Dot Product of vectors a and b
func (v Vector3) Dot(a Vector3) float32 {
    return (a.X * v.X) + (a.Y * v.Y) + (a.Z * v.Z)
}

// Cross will return the Cross Product of vectors v and a
func (v Vector3) Cross(a Vector3) Vector3 {
    return Vector3 {
        X: (v.Y * a. Z) - (v.Z * a.Y),
        Y: -((v.X * a.Z) - (v.Z * a.X)),
        Z: (v.X * a.Y) - (v.Y * a.X) }
}

// UnitVector will return the unit vector of a
func (v Vector3) UnitVector() Vector3 {
    len := float32(v.Length())
    return Vector3 {
        X: v.X / len,
        Y: v.Y / len,
        Z: v.Z / len }
}

// AsColor converts a vector to RGBA color values
func (v Vector3) AsColor() color.RGBA {
    return color.RGBA{uint8(v.X * math.MaxUint8), uint8(v.Y * math.MaxUint8), uint8(v.Z * math.MaxUint8), 255}
}