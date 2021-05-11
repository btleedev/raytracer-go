package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
)

type Light interface {
	hasPosition() bool
	getPosition() *r3.Vec
	getColorFrac() r3.Vec
	getLightIntensity() float64
	getSpecularLightIntensity() float64
	isPointVisible(point *r3.Vec, bvh *boundingVolumeHierarchy, monteCarloVariance *r3.Vec) bool
}

type AmbientLight struct {
	ColorFrac      r3.Vec
	LightIntensity float64
}

type PointLight struct {
	ColorFrac              r3.Vec
	Position               r3.Vec
	LightIntensity         float64
	SpecularLightIntensity float64
}

type SpotLight struct {
	ColorFrac              r3.Vec
	Position               r3.Vec
	LightIntensity         float64
	SpecularLightIntensity float64
	Direction              r3.Vec
	Angle                  float64 // specified in degrees
}

func (a AmbientLight) hasPosition() bool {
	return false
}

func (a AmbientLight) getPosition() *r3.Vec {
	return &r3.Vec{}
}

func (a AmbientLight) getColorFrac() r3.Vec {
	return a.ColorFrac
}

func (a AmbientLight) getLightIntensity() float64 {
	return a.LightIntensity
}

func (a AmbientLight) getSpecularLightIntensity() float64 {
	return 0
}

func (a AmbientLight) isPointVisible(point *r3.Vec, bvh *boundingVolumeHierarchy, monteCarloVariance *r3.Vec) bool {
	return true
}

func (p PointLight) hasPosition() bool {
	return true
}

func (p PointLight) getPosition() *r3.Vec {
	return &p.Position
}

func (p PointLight) getColorFrac() r3.Vec {
	return p.ColorFrac
}

func (p PointLight) getLightIntensity() float64 {
	return p.LightIntensity
}

func (p PointLight) getSpecularLightIntensity() float64 {
	return p.SpecularLightIntensity
}

func (p PointLight) isPointVisible(point *r3.Vec, bvh *boundingVolumeHierarchy, monteCarloVariance *r3.Vec) bool {
	shiftedPosition := r3.Add(p.Position, *monteCarloVariance)
	return doesReachLight(point, &shiftedPosition, bvh)
}

func (s SpotLight) hasPosition() bool {
	return true
}

func (s SpotLight) getPosition() *r3.Vec {
	return &s.Position
}

func (s SpotLight) getColorFrac() r3.Vec {
	return s.ColorFrac
}

func (s SpotLight) getLightIntensity() float64 {
	return s.LightIntensity
}

func (s SpotLight) getSpecularLightIntensity() float64 {
	return s.SpecularLightIntensity
}

func (s SpotLight) isPointVisible(point *r3.Vec, bvh *boundingVolumeHierarchy, monteCarloVariance *r3.Vec) bool {
	shiftedPosition := r3.Add(s.Position, *monteCarloVariance)
	reachesLight := doesReachLight(point, &shiftedPosition, bvh)

	// get angle between light direction vector and vector of light to point
	lightDirection := r3.Unit(s.Direction)
	lightPositionToShape := r3.Unit(r3.Sub(*point, shiftedPosition))
	angle := angleBetweenVectors(&lightDirection, &lightPositionToShape)
	return reachesLight && angle <= s.Angle
}

// unit is in degrees
func angleBetweenVectors(a, b *r3.Vec) float64 {
	aLengthSqrd := math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
	bLengthSqrd := math.Sqrt(b.X*b.X + b.Y*b.Y + b.Z*b.Z)
	angleRadians := math.Acos(r3.Dot(*a, *b) / (aLengthSqrd * bLengthSqrd))
	return angleRadians * 180 / math.Pi
}

func doesReachLight(origin *r3.Vec, lightPosition *r3.Vec, bvh *boundingVolumeHierarchy) bool {
	lightDirection := r3.Sub(*lightPosition, *origin)
	unitLightDirection := r3.Unit(lightDirection)
	r := ray{
		p:         *origin,
		direction: unitLightDirection,
	}
	hit, hitRecord := bvh.trace(
		&r,
		0.01, // don't let the shadow ray hit the same object
	)
	if !hit {
		return true
	}

	lengthFromPointToLight := r3.Dot(lightDirection, lightDirection)
	hitPointDirection := r3.Sub(hitRecord.p, *origin)
	lengthFromPointToHitObject := r3.Dot(hitPointDirection, hitPointDirection)
	return lengthFromPointToLight < lengthFromPointToHitObject
}
