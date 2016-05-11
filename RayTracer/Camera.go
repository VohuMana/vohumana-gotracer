package raytracer

import
(
    "encoding/json"
    "io"
    "math"
    "os"
)

// Camera is a struct to contain info about our virtual camera
type Camera struct {
    Origin, ImagePlaneHorizontal, ImagePlaneVertical, UpperLeftCorner Vector3
}

type cameraConfig struct {
    LookFrom, LookAt Vector3
    Fov float32
}

// GlobalCamera is the main camera object used in rendering
var GlobalCamera Camera

// cameraSettings is the global camera settings object
var cameraSettings cameraConfig

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
            X: 2.0 * halfWidth,
            Y: 0.0,
            Z: 0.0 },
        ImagePlaneVertical: Vector3 {
            X: 0.0,
            Y: -2.0 * halfHeight,
            Z: 0.0 },
        UpperLeftCorner: Vector3 {
            X: -halfWidth,
            Y: halfHeight,
            Z: -1.0 } }
}

// CreateCameraFromPos will create a camera looking at a point from another point
func CreateCameraFromPos(lookat, lookfrom, upVec Vector3, vFov, aspectRatio float32) Camera {
    cameraSettings.LookAt = lookat
    cameraSettings.LookFrom = lookfrom
    cameraSettings.Fov = vFov
    
    theta := ConvertDegreesToRadians(vFov)
    halfHeight := float32(math.Tan(float64(theta / 2.0)))
    halfWidth := aspectRatio * halfHeight
    w := lookat.Subtract(lookfrom).UnitVector()
    u := upVec.Cross(w).UnitVector()
    v := u.Cross(w).UnitVector()
    imageVert := v.Scale(2.0 * halfHeight)
    imageHoriz := u.Scale(-2.0 * halfWidth)
    corner := lookfrom.Subtract(v.Scale(halfHeight)).Subtract(u.Scale(-halfWidth)).Add(w)
    return Camera {
        Origin: lookfrom,
        ImagePlaneHorizontal: imageHoriz,
        ImagePlaneVertical: imageVert,
        UpperLeftCorner: corner }       
}

// ExportCamera will export the current global camera
func ExportCamera(filename string) {
    cameraString, err := json.Marshal(cameraSettings)
    checkError(err)
    
    cameraFile, err := os.Create(filename)
    checkError(err)
    defer cameraFile.Close()
    
    cameraFile.Write(cameraString)
}

// ImportCamera will import a camera json file into the global camera
func ImportCamera(filename string) {
    cameraFile, err := os.Open(filename)
    checkError(err)
    defer cameraFile.Close()
    
    info, err := cameraFile.Stat()
    checkError(err)
    
    contents := make([]byte, info.Size())
    
    _, err = cameraFile.Read(contents)
    if (err != nil && io.EOF != err) {
        checkError(err)   
    }

    err = json.Unmarshal(contents, &cameraSettings)
    checkError(err)
    
    GlobalCamera = CreateCameraFromPos(
        cameraSettings.LookAt, 
        cameraSettings.LookFrom,
        Vector3 {
            X: 0.0,
            Y: 1.0,
            Z: 0.0 },
       cameraSettings.Fov,
       float32(Settings.WidthInPixels) / float32(Settings.HeightInPixels))
}