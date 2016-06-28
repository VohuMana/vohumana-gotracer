package raytracer

import
(
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