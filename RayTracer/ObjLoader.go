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
    
    var verticies []Vector3
    var triangles []Triangle
    // var vertexNormals []Vector3
    
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
            fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d", &vertexIndex1, &textureIndex1, &normalIndex1, &vertexIndex2, &textureIndex2, &normalIndex2, &vertexIndex3, &textureIndex3, &normalIndex3)
            
            triangles = append(triangles, NewTriangle(
                verticies[vertexIndex1 - 1],
                verticies[vertexIndex2 - 1],
                verticies[vertexIndex3 - 1],
                Phong{}))
        }
    }
    
    return triangles
}