package raytracer

// Camera is a struct to contain info about our virtual camera
type Camera struct {
    Origin, ImagePlaneHorizontal, ImagePlaneVertical Vector3
}