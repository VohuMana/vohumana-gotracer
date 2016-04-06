package raytracer

// CollidableObject is an interface for objects that want to be able to collide with rays
type CollidableObject interface {
    TestIntersection(r Ray) (bool, IntersectionRecord)
}

// IntersectionRecord is an object that contains data about where a ray hit an object
type IntersectionRecord struct {
    T float32
    Point Vector3
    Normal Vector3
}