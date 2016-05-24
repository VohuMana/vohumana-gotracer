package raytracer

// LightList is a struct that holds a list of lights
type LightList struct {
    lights map[string]Light
}

func (l *LightList) addObject(name string, light Light) {
    if (nil == l.lights) {
        l.lights = make(map[string]Light)
    }
    
    l.lights[name] = light
}

// Light is a interface for light objects
type Light interface {
    GetColor() Vector3
    GetPosition() Vector3
    GetDirection() (Vector3, bool)
    GetPower() float32
    CalculateColor(r Ray, m Phong, i IntersectionRecord) Vector3
}