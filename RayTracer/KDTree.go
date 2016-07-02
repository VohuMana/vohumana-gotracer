package raytracer

import
(
    "image/color"
    "sort"
)

// KDTreeNode is the node that makes up the entire KDTree
type KDTreeNode struct {
    Left, Right *KDTreeNode
    Object []Triangle
    AABB BoundingBox
}

// KDTree is a structure that can be used for spatial partitioning
type KDTree struct {
    Root *KDTreeNode
}

func calculateTriangleMidpoint(t Triangle) Vector3 {
    return t.Position1.Add(t.Position2).Add(t.Position3).Scale(1.0 / 3.0)
}

// GenerateKDTree will generate a KDTree based on the passed in objects
func GenerateKDTree(tris []Triangle) KDTree {
    return KDTree {
        Root: buildKDTree(tris, 0) }
}

func buildKDTree(tris []Triangle, depth uint) *KDTreeNode {
    
    // If there are no tirangles in the incoming tris return an empty node
    if (len(tris) == 0) {
        return nil
    } else if (len(tris) <= 8) {
        // If the size is 1, create a node with 1 triangle
        return createNewKDNode(tris)
    }

    // Choose which axis to split on
    axis := depth % 3

    // Sort the triangles based on their midpoints
    switch axis {
        case 0:
            sort.Sort(ByMidpointX(tris))
            break

        case 1:
            sort.Sort(ByMidpointY(tris))
            break

        case 2:
            sort.Sort(ByMidpointZ(tris))
            break
    }

    // Send the first half of triangles to the left and the other half to the right
    midpoint := len(tris) / 2
    newNode := createNewKDNode(tris)
    newNode.Left = buildKDTree(tris[:midpoint], depth + 1)
    newNode.Right = buildKDTree(tris[midpoint:], depth + 1)
    return newNode
}

func createNewKDNode(tris []Triangle) *KDTreeNode {
    newNode := new(KDTreeNode)

    newNode.AABB = GenerateBoundingBoxFromTris(tris)
    newNode.Object = tris
    newNode.Left = nil
    newNode.Right = nil

    return newNode
}

// TestIntersection will check if anything has collided with the tree
func (k KDTree) TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    return testTreeIntersection(k.Root, &r, tMin, tMax)
}

func testTreeIntersection(treeNode *KDTreeNode, r *Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    // If the ray is colliding with the bounding box go deeper else no collision
    if (treeNode != nil && treeNode.AABB.IsRayColliding(*r)) {
        // Check if this node has children, if so continue traversal otherwise we are at a leaf node and need to check for collision with the triangles
        if (treeNode.Left != nil || treeNode.Right != nil) {
            hitLeft, recordLeft := testTreeIntersection(treeNode.Left, r, tMin, tMax)
            hitRight, recordRight := testTreeIntersection(treeNode.Right, r, tMin, tMax)

            if (hitLeft && hitRight) {
                if (recordLeft.T < recordRight.T) {
                    return true, recordLeft
                }

                return true, recordRight
            } else if (hitRight) {
                return true, recordRight
            }

            return hitLeft, recordLeft
        }
        
        collisionDetected := false
        var closestHitRecord IntersectionRecord
        closestT := tMax
        arrayLen := len(treeNode.Object)

        for i := 0; i < arrayLen; i++ {
            isColliding, hitRecord := treeNode.Object[i].TestIntersection(*r, tMin, closestT) 
            if (isColliding) {
                collisionDetected = true
                closestHitRecord = hitRecord
                closestT = hitRecord.T
            }
        }

        return collisionDetected, closestHitRecord
    }

    return false, IntersectionRecord{}
}

// GetColor will return the color of the triangle that intersected
func (k KDTree) GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA {
    return i.Object.GetColor(r, i, bounces)
}

// ByMidpointX is used to sort triangles based on their midpoint in the X axis
type ByMidpointX []Triangle

func (t ByMidpointX) Len() int { 
    return len(t) 
}

func (t ByMidpointX) Swap(i, j int) { 
    t[i], t[j] = t[j], t[i] 
}

func (t ByMidpointX) Less(i, j int) bool { 
    return t[i].Midpoint.X < t[j].Midpoint.X
}

// ByMidpointY is used to sort triangles based on their midpoint in the Y axis
type ByMidpointY []Triangle

func (t ByMidpointY) Len() int { 
    return len(t) 
}

func (t ByMidpointY) Swap(i, j int) { 
    t[i], t[j] = t[j], t[i] 
}

func (t ByMidpointY) Less(i, j int) bool { 
    return t[i].Midpoint.Y < t[j].Midpoint.Y
}

// ByMidpointZ is used to sort triangles based on their midpoint in the Z axis
type ByMidpointZ []Triangle

func (t ByMidpointZ) Len() int { 
    return len(t) 
}

func (t ByMidpointZ) Swap(i, j int) { 
    t[i], t[j] = t[j], t[i] 
}

func (t ByMidpointZ) Less(i, j int) bool { 
    return t[i].Midpoint.Z < t[j].Midpoint.Z
}