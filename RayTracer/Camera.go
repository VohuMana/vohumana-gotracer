package raytracer

import
(
    "math"
)

// Camera is a struct to contain info about our virtual camera
type Camera struct {
    Origin, ImagePlaneHorizontal, ImagePlaneVertical, UpperLeftCorner Vector3
}

// ConvertDegreesToRadians will convert the given degrees to radians
func ConvertDegreesToRadians(degrees float32) float32 {
    return degrees * float32(math.Pi / 180.0)
}

// CreateCamera will create a camera object with the vertical feild of view and aspect ratio requested
func CreateCamera(vFov, aspectRatio float32) Camera {
    theta := ConvertDegreesToRadians(vFov)
    halfHeight := float32(math.Tan(float64(theta / 2.0)))
    halfWidth := aspectRatio * halfHeight
    
    return Camera {
        Origin: Vector3 {
            X: 0.0,
            Y: 0.0,
            Z: 0.0 },
        ImagePlaneHorizontal: Vector3 {
            X: -2.0 * halfWidth,
            Y: 0.0,
            Z: 0.0 },
        ImagePlaneVertical: Vector3 {
            X: 0.0,
            Y: 2.0 * halfHeight,
            Z: 0.0 },
        UpperLeftCorner: Vector3 {
            X: halfWidth,
            Y: -halfHeight,
            Z: -1.0 } }
}

// CreateCameraFromPos will create a camera looking at a point from another point
func CreateCameraFromPos(lookat, lookfrom, upVec Vector3, vFov, aspectRatio float32) Camera {
    theta := ConvertDegreesToRadians(vFov)
    halfHeight := float32(math.Tan(float64(theta / 2.0)))
    halfWidth := aspectRatio * halfHeight
    w := lookfrom.Subtract(lookat).UnitVector()
    u := upVec.Cross(w).UnitVector()
    v := u.Cross(w).UnitVector()
    imageVert := v.Scale(2.0 * halfHeight)
    imageHoriz := u.Scale(2.0 * halfWidth)
    corner := lookfrom.Add(u.Scale(halfHeight)).Subtract(v.Scale(halfHeight)).Subtract(w)
    return Camera {
        Origin: lookfrom,
        ImagePlaneHorizontal: imageHoriz,
        ImagePlaneVertical: imageVert,
        UpperLeftCorner: corner }       
}