package raytracer

import
(
    "image/color"
)

// TriangleMesh is an object that contains a list of triangles that make up that mesh
type TriangleMesh struct {
    TriangleList []Triangle
}

// NewTriangleMesh creates a new triangle mesh from an obj file
func NewTriangleMesh(objFilename string, mat Material) TriangleMesh {
    meshTris := LoadObjFile(objFilename)

    // for _, tri := range meshTris {
    //     tri.Properties = mat
    // }

    for i := 0; i < len(meshTris); i++ {
        meshTris[i].Properties = mat
    }

    return NewTriangleMeshFromTriangleList(meshTris)
}

// NewTriangleMeshFromTriangleList creates a new triangle mesh from a list of triangles
func NewTriangleMeshFromTriangleList(triList []Triangle) TriangleMesh {
    return TriangleMesh {
        TriangleList: triList }
}

// ApplyMatrix3 will apply the given matrix to all triangles in the mesh
func (mesh TriangleMesh) ApplyMatrix3(mtx Matrix3) TriangleMesh {
    var newTris []Triangle
    
    for _, tri := range mesh.TriangleList {
        scaledTriangle := NewTriangle(
            mtx.MultiplyVector3(tri.Position1),
            mtx.MultiplyVector3(tri.Position2),
            mtx.MultiplyVector3(tri.Position3),
            tri.Properties)
        newTris = append(newTris, scaledTriangle)
    }

    return NewTriangleMeshFromTriangleList(newTris)
}

// Translate will move an object in 3D space, this should be done after scaling and rotating
func (mesh TriangleMesh) Translate(translation Vector3) TriangleMesh {
    var newTris []Triangle

    for _, tri := range mesh.TriangleList {
        translatedTriangle := NewTriangle(
            tri.Position1.Add(translation),
            tri.Position2.Add(translation),
            tri.Position3.Add(translation),
            tri.Properties)
        newTris = append(newTris, translatedTriangle)
    }

    return NewTriangleMeshFromTriangleList(newTris)
}

// TestIntersection will test if a ray is colliding with the triangle
func (mesh TriangleMesh) TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    collisionDetected := false
    var closestHitRecord IntersectionRecord
    closestT := tMax
    
    for _, tri := range mesh.TriangleList {
        isColliding, hitRecord := tri.TestIntersection(r, tMin, closestT) 
        if (isColliding) {
            collisionDetected = true
            closestHitRecord = hitRecord
            closestT = hitRecord.T
        }
    }
    
    return collisionDetected, closestHitRecord
}

// GetColor will get the color at the collision point
func (mesh TriangleMesh) GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA {
    return i.Object.GetColor(r, i, bounces)
}