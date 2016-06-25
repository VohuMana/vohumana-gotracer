package raytracer

// KDTreeNode is the node that makes up the entire KDTree
type KDTreeNode struct {
    Left, Right *KDTreeNode
    Object Triangle
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
    // If there are no tirangles in the incoming tris return an empty node
    // If the size is 1, create a node with 1 triangle

    // Choose which axis to split on
    // Sort the triangles based on their midpoints
    // Send the first half of triangles to the left and the other half to the right

    return KDTree {
        Root: nil }
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