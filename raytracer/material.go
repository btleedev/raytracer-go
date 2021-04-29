package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
	"math/rand"
)

type material interface {
	scatter(r *ray, hitRecord *hitRecord, shapes *[]shape, lights *[]light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec)
}

type diffuse struct {
	albedo r3.Vec
	color  r3.Vec
}

type metal struct {
	albedo r3.Vec
	fuzz   float64
}

type dielectric struct {
	refractiveIndex float64
}

type phongBlinn struct {
	specValue     float64
	specShininess float64
	color         r3.Vec
}

func (d diffuse) scatter(r *ray, hitRecord *hitRecord, shapes *[]shape, lights *[]light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec) {
	target := r3.Add(hitRecord.p, r3.Add(hitRecord.normal, randomInUnitSphere()))
	// set shouldTrace = true to allow diffuse materials to scatter as well, turned off to compute shadows
	return false, d.albedo, ray{p: hitRecord.p, direction: r3.Sub(target, hitRecord.p)}, d.color
}

func (m metal) scatter(r *ray, hitRecord *hitRecord, shapes *[]shape, lights *[]light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec) {
	correctedFuzz := 1.0
	if m.fuzz < 1.0 {
		correctedFuzz = m.fuzz
	}
	directionNormalized := r3.Unit(r.direction)
	reflectedRay := reflected(&directionNormalized, &hitRecord.normal)
	return r3.Dot(reflectedRay, hitRecord.normal) > 0, m.albedo, ray{p: hitRecord.p, direction: r3.Add(reflectedRay, r3.Scale(correctedFuzz, randomInUnitSphere()))}, r3.Vec{}
}

func (d dielectric) scatter(r *ray, hitRecord *hitRecord, shapes *[]shape, lights *[]light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec) {
	outwardNormal := r3.Vec{}
	niOverNt := 0.0
	reflectProb := 0.0
	cosine := 0.0
	reflectedVec := reflected(&r.direction, &hitRecord.normal)
	if r3.Dot(r.direction, hitRecord.normal) > 0 {
		outwardNormal = r3.Scale(-1, hitRecord.normal)
		niOverNt = d.refractiveIndex
		cosine = d.refractiveIndex * r3.Dot(r3.Unit(r.direction), hitRecord.normal) / math.Sqrt(r.direction.X*r.direction.X+r.direction.Y*r.direction.Y+r.direction.Z*r.direction.Z)
	} else {
		outwardNormal = hitRecord.normal
		niOverNt = 1.0 / d.refractiveIndex
		cosine = -1 * r3.Dot(r3.Unit(r.direction), hitRecord.normal) / math.Sqrt(r.direction.X*r.direction.X+r.direction.Y*r.direction.Y+r.direction.Z*r.direction.Z)
	}

	shouldRefract, refractedVec := refracted(&r.direction, &outwardNormal, niOverNt)
	if shouldRefract {
		reflectProb = schlick(cosine, d.refractiveIndex)
	} else {
		reflectProb = 1.0
	}

	if rand.Float64() < reflectProb {
		return true, r3.Vec{X: 1.0, Y: 1.0, Z: 1.0}, ray{p: r3.Add(hitRecord.p, r3.Scale(0.00001, reflectedVec)), direction: reflectedVec}, r3.Vec{}
	} else {
		return true, r3.Vec{X: 1.0, Y: 1.0, Z: 1.0}, ray{p: r3.Add(hitRecord.p, r3.Scale(0.00001, refractedVec)), direction: refractedVec}, r3.Vec{}
	}
}

func (p phongBlinn) scatter(r *ray, hitRecord *hitRecord, shapes *[]shape, lights *[]light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec) {
	c := r3.Vec{}
	for _, light := range *lights {
		if light.hasPosition() {
			monteCarloRepetitions := softShadowMonteCarloRepetitions
			monteCarloMaxLength := softShadowMonteCarloMaxLengthDeviation
			for i := 0; i < monteCarloRepetitions; i++ {
				hitPoint := hitRecord.p
				monteCarloVariance := r3.Scale(monteCarloMaxLength, randomInUnitSphere())
				if light.isPointVisible(&hitPoint, shapes, &monteCarloVariance) {
					lightPosition := *light.getPosition()
					lightToPoint := r3.Sub(lightPosition, hitPoint)
					lightDirection := r3.Unit(lightToPoint)
					lightDistanceSqrd := lightToPoint.X*lightToPoint.X + lightToPoint.Y*lightToPoint.Y + lightToPoint.Z*lightToPoint.Z
					viewDirection := r3.Unit(r3.Sub(r.p, hitPoint))
					blinnDirection := r3.Unit(r3.Add(lightDirection, viewDirection))
					blinnTerm := math.Max(r3.Dot(r3.Unit(hitRecord.normal), blinnDirection), 0.0)
					phongTerm := light.getLightIntensity() * p.specValue * math.Pow(blinnTerm, p.specShininess) / lightDistanceSqrd
					lambertTerm := light.getLightIntensity() * math.Max(0.0, r3.Dot(lightDirection, r3.Unit(hitRecord.normal))) / lightDistanceSqrd

					c = r3.Add(c, r3.Scale(1/float64(monteCarloRepetitions), r3.Scale(phongTerm, light.getColorFrac())))
					c = r3.Add(c, r3.Scale(1/float64(monteCarloRepetitions), r3.Scale(lambertTerm, p.color)))
				}
			}
		} else {
			// treat as ambient light - no position so we assume it can reach us
			c = r3.Add(c, r3.Scale(light.getLightIntensity(), light.getColorFrac()))
		}
	}
	c.X = math.Min(1.0, c.X)
	c.Y = math.Min(1.0, c.Y)
	c.Z = math.Min(1.0, c.Z)
	return false, r3.Vec{}, ray{}, c
}

func randomInUnitSphere() r3.Vec {
	p := r3.Vec{}
	for {
		p = r3.Sub(r3.Scale(2, r3.Vec{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()}), r3.Vec{X: 1, Y: 1, Z: 1})
		if p.X*p.X+p.Y*p.Y+p.Z*p.Z < 1.0 {
			break
		}
	}
	return p
}

func reflected(v *r3.Vec, n *r3.Vec) r3.Vec {
	return r3.Sub(*v, r3.Scale(2*r3.Dot(*v, *n), *n))
}

func refracted(v *r3.Vec, n *r3.Vec, niOverNt float64) (b bool, refracted r3.Vec) {
	uv := r3.Unit(*v)
	dt := r3.Dot(uv, *n)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		return true, r3.Sub(r3.Scale(niOverNt, r3.Sub(uv, r3.Scale(dt, *n))), r3.Scale(math.Sqrt(discriminant), *n))
	} else {
		return false, r3.Vec{}
	}
}

func schlick(cosine float64, refractiveIndex float64) float64 {
	r0 := (1 - refractiveIndex) / (1 + refractiveIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
