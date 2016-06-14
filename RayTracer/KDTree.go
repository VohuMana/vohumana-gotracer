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

// GenerateKDTree will generate a KDTree based on the passed in objects
func GenerateKDTree(tris []Triangle) KDTree {

    return KDTree {
        Root: nil }
}
