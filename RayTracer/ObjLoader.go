package raytracer

import
(
    "bufio"
    "fmt"
    "os"
)

func LoadObjFile(filename string) []Triangle {
    file, err := os.Open(filename)
    checkError(err)
    defer file.Close()
    
    reader := bufio.NewReader(file)
    scanner := bufio.NewScanner(reader)
    
    verts := false
    normals := false
    textures := false

    var verticies []Vector3
    var triangles []Triangle
    var vertexNormals []Vector3
    var textureCoords []Vector3
    
    for scanner.Scan() {
        line := scanner.Text()

        var text string
        fmt.Sscanf(line, "%s", &text)
        if (text == "#") {
            continue
        } else if (text == "v") {
            var X float32
            var Y float32
            var Z float32
            fmt.Sscanf(line, "%s %f %f %f", &text, &X, &Y, &Z)
            verticies = append(verticies, NewVector3(X,Y,Z))
            verts = true
        } else if (text == "vn") {
            var X float32
            var Y float32
            var Z float32
            fmt.Sscanf(line, "%s %f %f %f", &text, &X, &Y, &Z)
            vertexNormals = append(vertexNormals, NewVector3(X,Y,Z))
            normals = true
        } else if (text == "vt") {
            var X float32
            var Y float32
            var Z float32
            fmt.Sscanf(line, "%s %f %f %f", &text, &X, &Y, &Z)
            textureCoords = append(textureCoords, NewVector3(X,Y,Z))
            textures = true
        } else if (text == "f") {
            var vertexIndex1 int
            var textureIndex1 int
            var normalIndex1 int
            
            var vertexIndex2 int
            var textureIndex2 int
            var normalIndex2 int
            
            var vertexIndex3 int
            var textureIndex3 int
            var normalIndex3 int

            if (verts && !normals && !textures) {
                fmt.Sscanf(line, "f %d %d %d", &vertexIndex1, &vertexIndex2, &vertexIndex3)
            } else if (verts && normals && !textures) {
                fmt.Sscanf(line, "f %d/%d %d/%d %d/%d", &vertexIndex1, &normalIndex1, &vertexIndex2, &normalIndex2, &vertexIndex3, &normalIndex3)    
            } else {
                fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d", &vertexIndex1, &textureIndex1, &normalIndex1, &vertexIndex2, &textureIndex2, &normalIndex2, &vertexIndex3, &textureIndex3, &normalIndex3)
            }

            if (normals) {
                normalAverage := vertexNormals[normalIndex1 - 1].Add(vertexNormals[normalIndex2 - 1]).Add(vertexNormals[normalIndex3 - 1]).Scale(1.0 / 3.0)

                triangles = append(triangles, NewTriangleWithNormal(
                    verticies[vertexIndex1 - 1],
                    verticies[vertexIndex2 - 1],
                    verticies[vertexIndex3 - 1],
                    normalAverage,
                    Phong{}))
            } else {
                triangles = append(triangles, NewTriangle(
                    verticies[vertexIndex1 - 1],
                    verticies[vertexIndex2 - 1],
                    verticies[vertexIndex3 - 1],
                    Phong{}))
            }
            
            
        }
    }
    
    return triangles
}