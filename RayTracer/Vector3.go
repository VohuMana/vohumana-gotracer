package RayTracer

import
(
    "math"
)

// Vector3 is a basic 3D vector type with members x, y, and z
type Vector3 struct {
    x, y, z float32
}

// Add adds vectors a and b like a + b
func (v *Vector3) Add(a, b *Vector3) *Vector3 {
    v.x = a.x + b.x
    v.y = a.y + b.y
    v.z = a.z + b.z
    return v
}

// Subtract subtracts vectors a and b like a - b
func (v *Vector3) Subtract(a, b *Vector3) *Vector3 {
    v.x = a.x - b.x
    v.y = a.y - b.y
    v.z = a.z - b.z
    return v
}

// Multiply multiplies vectors a and b like a * b
func (v *Vector3) Multiply(a, b *Vector3) *Vector3  {
    v.x = a.x * b.x
    v.y = a.y * b.y
    v.z = a.z * b.z
    return v
}

// Divide divides vectors a and b like a / b
func (v *Vector3) Divide(a, b *Vector3) *Vector3 {
    v.x = a.x / b.x
    v.y = a.y / b.y
    v.z = a.z / b.z
    return v
}

// Length will get the length of the vector as defined by sqrt(x^2 + y^2 + z^2)
func (v *Vector3) Length(a *Vector3) float64 {
    return math.Sqrt(float64(a.x * a.x) + float64(a.y * a.y) + float64(a.z * a.z))
}

// SquareLength will return the length of the vector squared.  x^2 + y^2 + z^2
func (v *Vector3) SquareLength(a *Vector3) float64 {
    return float64(a.x * a.x) + float64(a.y * a.y) + float64(a.z * a.z)
}

// Dot will return the Dot Product of vectors a and b
func (v *Vector3) Dot(a, b *Vector3) float32 {
    return (a.x * b.x) + (a.y * b.y) + (a.z * b.z)
}

// Cross will return the Cross Product of vectors a and b
func (v *Vector3) Cross(a, b *Vector3) *Vector3 {
    v.x = (a.y * b. z) - (a.z * b.y)
    v.y = -((a.x * b.z) - (a.z * b.x))
    v.z = (a.x * b.y) - (a.y * b.x)
    return v
}

// UnitVector will return the unit vector of a
func (v *Vector3) UnitVector(a *Vector3) *Vector3 {
    var len float32
    len = float32(a.Length(a))
    v.x = a.x / len
    v.y = a.y / len
    v.z = a.z / len
    return v
}