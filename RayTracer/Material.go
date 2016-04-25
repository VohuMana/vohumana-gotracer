package raytracer

import
(
    "image/color"
    "math"
    "math/rand"
)

// Material is a interface that returns where a scattered ray will be when it reflects off an object
type Material interface {
    Scatter(r Ray, i IntersectionRecord) Ray
    GetAttenuation() Vector3
}

// Lambertian is a type of material that scatters rays randomly, used for diffuse objectss
type Lambertian struct {
    Color color.RGBA
    Attenuation Vector3
}

// Metal is a type of material that reflects rays
type Metal struct {
    Color color.RGBA
    Fuzziness float32
    Attenuation Vector3
}

// Dielectric is a type of material that refracts rays
type Dielectric struct {
    RefractiveIndex float32
    Attenuation Vector3
}

// calculateReflectionVector calulates a reflection vector with direction d and normal n. r = d - 2(n dot d)*n
func calculateReflectionVector(d, n Vector3) Vector3 {
    return d.Subtract(n.Scale(2.0 * d.Dot(n))).UnitVector()
}

// getNormalAsColor give a normal n, return the color value for that normal
func getNormalAsColor(n Vector3) color.RGBA {
    // Render normals
    c := n.Scale(0.5)
    c = c.Add(Vector3{X:1.0, Y:1.0, Z:1.0})
    
    c.X = float32(math.Min(1.0, float64(c.X)))
    c.Y = float32(math.Min(1.0, float64(c.Y)))
    c.Z = float32(math.Min(1.0, float64(c.Z)))
    return c.AsColor()
}

func createUnitSphereVector() Vector3 {
    return Vector3 {
        rand.Float32(),
        rand.Float32(),
        rand.Float32() }.Subtract(Vector3 {
            1.0,
            1.0,
            1.0 }).Multiply(Vector3 {
                2.0,
                2.0,
                2.0 })
}

func  randomVectorInUnitSphere() Vector3 {
    p := createUnitSphereVector()
    
    for p.Dot(p) >= 1.0 {
        p = createUnitSphereVector()
    }
    
    return p
}

func restrictValues(num, min, max float32) float32 {
    num = float32(math.Min(float64(num), float64(max)))
    num = float32(math.Max(float64(num), float64(min)))
    return num
}

func calculateReflectionRay(r Ray, i IntersectionRecord, fuzziness float32) Ray {
    return Ray {
        Origin: i.Point,
        Direction: calculateReflectionVector(r.Direction, i.Normal).Add(randomVectorInUnitSphere().Scale(fuzziness)).UnitVector() }
}

func calculateDiffuseRay(i IntersectionRecord) Ray {
        target := i.Point.Add(i.Normal).Add(randomVectorInUnitSphere())
        return Ray {
            Origin: i.Point,
            Direction: target.Subtract(i.Point).UnitVector() }
}

func schlickReflectanceProbability(cosine, refractiveIndex float32) float32 {
    r0 := (1.0 - refractiveIndex) / (1.0 + refractiveIndex)
    r0 = r0 * r0
    return r0 + (1.0 - r0) * float32(math.Pow(float64(1.0 - cosine), 5))
}

func calculateRefractionVector(incomingVector, normal Vector3, niOverNt float32) (bool, Vector3) {
    dt := incomingVector.Dot(normal)
    discriminant := 1.0 - niOverNt * niOverNt * (1 - dt * dt)
    if (discriminant > 0.0) {
        // niOverNt * (incomingVector - normal * dt) - normal * sqrt(discriminant)
        return true, incomingVector.Subtract(normal.Scale(dt)).Scale(niOverNt).Subtract(normal.Scale(float32(math.Sqrt(float64(discriminant)))))
    }
    
    return false, Vector3{}
}

func calculateRefractedRay(r Ray, i IntersectionRecord, refractiveIndex float32) Ray {
    var outwardNormal Vector3
    var niOverNt float32
    var refractedRay Ray
    var cosine float32
    var reflectionProbability float32
    reflectionVector := calculateReflectionVector(r.Direction, i.Normal)
    
    if (r.Direction.Dot(i.Normal) > 0.0) {
        outwardNormal = i.Normal.Scale(-1.0)
        niOverNt = refractiveIndex
        cosine = refractiveIndex * r.Direction.Dot(i.Normal)
    } else {
        outwardNormal = i.Normal
        niOverNt = 1.0 / refractiveIndex
        cosine = -r.Direction.Dot(i.Normal)
    }
    
    isRefracted, refractedVector := calculateRefractionVector(r.Direction, outwardNormal, niOverNt)
    
    refractedRay.Origin = i.Point
    
    if (isRefracted) {
        reflectionProbability = schlickReflectanceProbability(cosine, refractiveIndex)
    } else {
        reflectionProbability = 1.0
    }
    
    if (rand.Float32() < reflectionProbability) {
        refractedRay.Direction = reflectionVector
    } else {
        refractedRay.Direction = refractedVector.UnitVector()
    }
    
    return refractedRay
}

// Scatter for lambertian materials
func (l Lambertian) Scatter(r Ray, i IntersectionRecord) Ray {
    return calculateDiffuseRay(i)
}

// GetAttenuation returns the diffuse attenuation
func (l Lambertian) GetAttenuation() Vector3 {
    return l.Attenuation
}

// Scatter for metal materials
func (m Metal) Scatter(r Ray, i IntersectionRecord) Ray {
    return calculateReflectionRay(r, i, m.Fuzziness)
}

// GetAttenuation gets the metal attenuation
func (m Metal) GetAttenuation() Vector3 {
    return m.Attenuation
}

// Scatter refracts rays for dielectric materials
func (d Dielectric) Scatter(r Ray, i IntersectionRecord) Ray {
    return calculateRefractedRay(r, i, d.RefractiveIndex)
}

// GetAttenuation gets the dielectric attenuation
func (d Dielectric) GetAttenuation() Vector3 {
    return d.Attenuation
}