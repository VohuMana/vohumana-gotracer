package raytracer

// CollisionList is a struct used to abstract the collision test loop
type CollisionList struct {
    collisionList map[string]CollidableObject
}

// AddObject adds a collidableobject to the collision map
func (c *CollisionList) AddObject(name string, obj CollidableObject) {
    if (nil == c.collisionList) {
        c.collisionList = make(map[string]CollidableObject)
    }
    
    c.collisionList[name] = obj
}

// TestCollision loops through all the objects testing if the ray is colliding with any and returns the nearest object
func (c CollisionList) TestCollision(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    collisionDetected := false
    var closestHitRecord IntersectionRecord
    closestT := tMax
    
    for _, obj := range c.collisionList {
        isColliding, hitRecord := obj.TestIntersection(r, tMin, closestT) 
        if (isColliding) {
            collisionDetected = true
            closestHitRecord = hitRecord
            closestT = hitRecord.T
        }
    }
    
    return collisionDetected, closestHitRecord
}