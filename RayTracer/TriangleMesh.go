package raytracer

import
(
    "image/color"
)

// TriangleMesh is an object that contains a list of triangles that make up that mesh
type TriangleMesh struct {
    TriangleList []Triangle
}

// ExportableTriangleMesh is a way to export a triangle mesh with out having to export all the triangles
type ExportableTriangleMesh struct {
    ObjectFileName string
    Scale, Translation Vector3
    Properties Material
}

// NewExportableTriangleMesh creates a new exportable triangle mesh
func NewExportableTriangleMesh(objFilePath string, scaleValue, translationValue Vector3, mat Material) ExportableTriangleMesh {
    return ExportableTriangleMesh {
        ObjectFileName: objFilePath,
        Scale: scaleValue,
        Translation: translationValue,
        Properties: mat }
}

// TestIntersection does nothing for an ExportableTriangleMesh
func (mesh ExportableTriangleMesh) TestIntersection(r Ray, tMin, tMax float32) (bool, IntersectionRecord) {
    return false, IntersectionRecord{}
}

// GetColor does nothing for an ExportableTriangleMesh
func (mesh ExportableTriangleMesh) GetColor(r Ray, i IntersectionRecord, bounces uint32) color.RGBA {
    return color.RGBA {R: 0, G:0, B:0, A:255}
}

func deserializeTriMesh(object map[string]interface{}) (TriangleMesh, bool) {
    var meshExport ExportableTriangleMesh
    var mesh TriangleMesh
    validMesh := true
    
    for name, object := range object {
        switch name {
            case "Scale":
                scale, ok := object.(map[string]interface{})
                if (true == ok) {
                    meshExport.Scale, ok = deserializeVector3(scale)
                    if (false == ok) {
                        validMesh = false
                        break;
                    }
                }
                
            case "Translation":
                translation, ok := object.(map[string]interface{})
                if (true == ok) {
                    meshExport.Translation, ok = deserializeVector3(translation)
                    if (false == ok) {
                        validMesh = false
                        break;
                    }
                }

            case "ObjectFileName":
                filename, ok := object.(string)
                if (true == ok) {
                    meshExport.ObjectFileName = string(filename)
                } else {
                    validMesh = false
                    break;
                }
                
            case "Properties":
                prop, ok := object.(map[string]interface{})
                if (true == ok) {
                    mat, ok := deserializeMaterial(prop)
                    if (true == ok) {
                        meshExport.Properties = mat
                    } else {
                        validMesh = false
                        break;
                    }    
                }
            
            default:
                validMesh = false
                break
        }
    }

    if (validMesh) {
        mesh = NewTriangleMesh(meshExport.ObjectFileName, meshExport.Properties)
        scaleMtx := ScaleMatrix(meshExport.Scale.X, meshExport.Scale.Y, meshExport.Scale.Z)
        mesh = mesh.ApplyMatrix3(scaleMtx)
        mesh = mesh.Translate(meshExport.Translation)
    }
    
    return mesh, validMesh
}

// NewTriangleMesh creates a new triangle mesh from an obj file
func NewTriangleMesh(objFilename string, mat Material) TriangleMesh {
    meshTris := LoadObjFile(objFilename)

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