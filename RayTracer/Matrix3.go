package raytracer

// Matrix3 is a 3x3 matrix
type Matrix3 struct {
    Row1, Row2, Row3 Vector3
}

// NewMatrix3 will generate a new 3x3 matrix with the provided rows
func NewMatrix3(row1, row2, row3 Vector3) Matrix3 {
    return Matrix3 {
        Row1: row1,
        Row2: row2,
        Row3: row3 }
}

// IdentityMatrix3 will generate a new 3x3 identity matrix
func IdentityMatrix3() Matrix3 {
    return NewMatrix3(
        NewVector3(1, 0, 0),
        NewVector3(0, 1, 0),
        NewVector3(0, 0, 1))
}

// MultiplyMatrix will multiply 2 matricies together m * a
func (m Matrix3) MultiplyMatrix(a Matrix3) Matrix3 {
    newRow1 := NewVector3(
        rowMultColumn(m.Row1, a.Row1.X, a.Row2.X, a.Row3.X),
        rowMultColumn(m.Row1, a.Row1.Y, a.Row2.Y, a.Row3.Y),
        rowMultColumn(m.Row1, a.Row1.Z, a.Row2.Z, a.Row3.Z))
    newRow2 := NewVector3(
        rowMultColumn(m.Row2, a.Row1.X, a.Row2.X, a.Row3.X),
        rowMultColumn(m.Row2, a.Row1.Y, a.Row2.Y, a.Row3.Y),
        rowMultColumn(m.Row2, a.Row1.Z, a.Row2.Z, a.Row3.Z))
    newRow3 := NewVector3(
        rowMultColumn(m.Row3, a.Row1.X, a.Row2.X, a.Row3.X),
        rowMultColumn(m.Row3, a.Row1.Y, a.Row2.Y, a.Row3.Y),
        rowMultColumn(m.Row3, a.Row1.Z, a.Row2.Z, a.Row3.Z))
    return NewMatrix3(newRow1, newRow2, newRow3)
}

// MultiplyVector3 will multiply a matirx and a vector
func (m Matrix3) MultiplyVector3(vec Vector3) Vector3 {
    return NewVector3(
        vec.X * m.Row1.X + vec.Y * m.Row1.Y + vec.Z * m.Row1.Z,
        vec.X * m.Row2.X + vec.Y * m.Row2.Y + vec.Z * m.Row2.Z,
        vec.X * m.Row3.X + vec.Y * m.Row2.Y + vec.Z * m.Row3.Z)
}

func rowMultColumn(row Vector3, x, y, z float32) float32 {
    return row.X * x + row.Y * y + row.Z * z
}