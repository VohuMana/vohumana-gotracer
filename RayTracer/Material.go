package raytracer

import
(
    "encoding/json"
    "image/color"
    "math"
    "math/rand"
)

// Material is a interface that returns where a scattered ray will be when it reflects off an object
type Material interface {
    GetColor(r Ray, i IntersectionRecord, w World, bounceDepth uint32) color.RGBA
}

// Lambertian is a type of material that scatters rays randomly, used for diffuse objectss
type Lambertian struct {
    Color color.RGBA
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

func calculatePhongLighting(p Phong, i IntersectionRecord, reflectRay Ray, w World, bounceDepth uint32) color.RGBA {
    var color Vector3
    
    // Add up the contribution of all light sources
    for _, light := range w.Lights.lights {
        color = color.Add(computeColorContributionFromLight(light, &i, &w, &p, &reflectRay))
    } 
    
    // Get the color of anything that might have been reflected
    reflectedColor := AsVector3(ShootRay(reflectRay, w, bounceDepth + 1)).Scale(p.Reflectivity)
    
    return color.Add(reflectedColor).AsColor()
}

func computeColorContributionFromLight(light Light, i *IntersectionRecord, w *World, p *Phong, reflectRay *Ray) Vector3 {
    var color Vector3
    
    for samples := uint32(0); samples < Settings.MaxLightRays; samples++ {
        lightVector := light.GetPosition().Subtract(i.Point)
        distance := float32(lightVector.Length())
        
        if samples == 0 {
            lightVector = lightVector.UnitVector()    
        } else {
            lightVector = lightVector.Add(randomVectorInUnitSphere().Scale(distance)).UnitVector()    
        }
        
        lightRay := Ray {
            Origin: i.Point,
            Direction: lightVector }
        
        isShadow, _ := w.TestCollision(lightRay, 0.001, distance)
        
        // Contribute no color if the object is in shadow
        if (isShadow) {
            // TODO: Add ambient
            continue
        }
        
        // Calculate the diffuse coefficient 
        diffuse := (1.0 - p.Reflectivity) * i.Normal.Dot(lightVector) * light.GetPower()
        
        // Calculate the vector to the camera
        cameraVector := GlobalCamera.Origin.Subtract(i.Point).UnitVector()
        
        foo := math.Pow(float64(reflectRay.Direction.Dot(cameraVector)), float64(p.Shininess))
        
        // Calculate the specular coefficient        
        specular := float32(math.Max(0.0, foo)) * light.GetPower()
        
        // Phong lighting for given point
        colorContribution := p.DiffuseColor.Scale(diffuse)
        colorContribution = colorContribution.Add(colorContribution.Scale(specular))
        colorContribution = colorContribution.Multiply(light.GetColor())
        
        // Add the color contribution from this sample to the overall color
        color = color.Add(colorContribution)
    }
    
    // Average the colors
    return color.Scale( 1.0 / float32(Settings.MaxLightRays))
}

// Scatter for lambertian materials
func (l Lambertian) Scatter(r Ray, i IntersectionRecord) Ray {
    return calculateDiffuseRay(i)
}

// GetAttenuation returns the diffuse attenuation
func (l Lambertian) GetAttenuation() Vector3 {
    return l.Attenuation
}

// GetEmission returns the emissive component for lights, lambertian has no emissive component
func (l Lambertian) GetEmission() Vector3 {
    return NewVector3(0.0, 0.0, 0.0)
}

// IsEmissive returns if the object is emissive or not
func (l Lambertian) IsEmissive() bool {
    return false
}

func deserializeMaterial(object map[string]interface{}) (Material, bool) {
    b, err := json.Marshal(object)
    if (err != nil) {
        checkError(err)
    }
    
    // If Fuzziness exists in the object then it must be a metal
    if nil != object["Fuzziness"] {
        var metal Metal
        if err := json.Unmarshal(b, &metal); err == nil {
            return metal, true
        }
    // If RefractiveIndex is in the object then it must be a dielectric
    } else if nil != object["RefractiveIndex"] {
        var dielectric Dielectric
        if err := json.Unmarshal(b, &dielectric); err == nil {
            return dielectric, true
        }
    // If Emission is in the object then it must be a dielectric
    } else if nil != object["Emission"] {
        var emissive Emissive
        if err := json.Unmarshal(b, &emissive); err == nil {
            return emissive, true
        }
    } else if nil != object["DiffuseColor"] {
        var phong Phong
        if err := json.Unmarshal(b, &phong); err == nil {
            return phong, true
        }
    }   
    
    return nil, false
}