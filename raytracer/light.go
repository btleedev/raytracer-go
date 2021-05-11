package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
)

type light interface {
	hasPosition() bool
	getPosition() *r3.Vec
	getColorFrac() r3.Vec
	getLightIntensity() float64
	getSpecularLightIntensity() float64
	isPointVisible(point *r3.Vec, bvh *boundingVolumeHierarchy, monteCarloVariance *r3.Vec) bool
}

type ambientLight struct {
	colorFrac      r3.Vec
	lightIntensity float64
}

type pointLight struct {
	colorFrac              r3.Vec
	position               r3.Vec
	lightIntensity         float64
	specularLightIntensity float64
}

type spotLight struct {
	colorFrac              r3.Vec
	position               r3.Vec
	lightIntensity         float64
	specularLightIntensity float64
	direction              r3.Vec
	angle                  float64 // specified in degrees
}

func (a ambientLight) hasPosition() bool {
	return false
}

func (a ambientLight) getPosition() *r3.Vec {
	return &r3.Vec{}
}

func (a ambientLight) getColorFrac() r3.Vec {
	return a.colorFrac
}

func (a ambientLight) getLightIntensity() float64 {
	return a.lightIntensity
}

func (a ambientLight) getSpecularLightIntensity() float64 {
	return 0
}

func (a ambientLight) isPointVisible(point *r3.Vec, bvh *boundingVolumeHierarchy, monteCarloVariance *r3.Vec) bool {
	return true
}

func (p pointLight) hasPosition() bool {
	return true
}

func (p pointLight) getPosition() *r3.Vec {
	return &p.position
}

func (p pointLight) getColorFrac() r3.Vec {
	return p.colorFrac
}

func (p pointLight) getLightIntensity() float64 {
	return p.lightIntensity
}

func (p pointLight) getSpecularLightIntensity() float64 {
	return p.specularLightIntensity
}

func (p pointLight) isPointVisible(point *r3.Vec, bvh *boundingVolumeHierarchy, monteCarloVariance *r3.Vec) bool {
	shiftedPosition := r3.Add(p.position, *monteCarloVariance)
	return doesReachLight(point, &shiftedPosition, bvh)
}

func (s spotLight) hasPosition() bool {
	return true
}

func (s spotLight) getPosition() *r3.Vec {
	return &s.position
}

func (s spotLight) getColorFrac() r3.Vec {
	return s.colorFrac
}

func (s spotLight) getLightIntensity() float64 {
	return s.lightIntensity
}

func (s spotLight) getSpecularLightIntensity() float64 {
	return s.specularLightIntensity
}

func (s spotLight) isPointVisible(point *r3.Vec, bvh *boundingVolumeHierarchy, monteCarloVariance *r3.Vec) bool {
	shiftedPosition := r3.Add(s.position, *monteCarloVariance)
	reachesLight := doesReachLight(point, &shiftedPosition, bvh)

	// get angle between light direction vector and vector of light to point
	lightDirection := r3.Unit(s.direction)
	lightPositionToShape := r3.Unit(r3.Sub(*point, shiftedPosition))
	angle := angleBetweenVectors(&lightDirection, &lightPositionToShape)
	return reachesLight && angle <= s.angle
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
