package raytracer

// CollidableObject is an interface for objects that want to be able to collide with rays
type CollidableObject interface {
    TestIntersection(r Ray) bool
}