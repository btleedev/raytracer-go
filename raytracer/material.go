package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
	"math/rand"
)

type Material interface {
	scatter(is *ImageSpec, r *ray, hitRecord *hitRecord, bvh *boundingVolumeHierarchy, lights *[]Light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec)
}

type Standard struct {
	Color r3.Vec
}

type Metal struct {
	Albedo r3.Vec
	Fuzz   float64
}

type Dielectric struct {
	RefractiveIndex float64
}

type PhongBlinn struct {
	Color         r3.Vec
	SpecularColor r3.Vec
	SpecHardness  float64
}

func (d Standard) scatter(is *ImageSpec, r *ray, hitRecord *hitRecord, bvh *boundingVolumeHierarchy, lights *[]Light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec) {
	return false, r3.Vec{}, ray{p: hitRecord.p, direction: r3.Vec{}}, d.Color
}

func (m Metal) scatter(is *ImageSpec, r *ray, hitRecord *hitRecord, bvh *boundingVolumeHierarchy, lights *[]Light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec) {
	correctedFuzz := 1.0
	if m.Fuzz < 1.0 {
		correctedFuzz = m.Fuzz
	}
	directionNormalized := r3.Unit(r.direction)
	reflectedRay := reflected(&directionNormalized, &hitRecord.normal)
	return r3.Dot(reflectedRay, hitRecord.normal) > 0, m.Albedo, ray{p: hitRecord.p, direction: r3.Add(reflectedRay, r3.Scale(correctedFuzz, randomInUnitSphere()))}, r3.Vec{}
}

func (d Dielectric) scatter(is *ImageSpec, r *ray, hitRecord *hitRecord, bvh *boundingVolumeHierarchy, lights *[]Light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec) {
	outwardNormal := r3.Vec{}
	niOverNt := 0.0
	reflectProb := 0.0
	cosine := 0.0
	reflectedVec := reflected(&r.direction, &hitRecord.normal)
	if r3.Dot(r.direction, hitRecord.normal) > 0 {
		outwardNormal = r3.Scale(-1, hitRecord.normal)
		niOverNt = d.RefractiveIndex
		cosine = d.RefractiveIndex * r3.Dot(r3.Unit(r.direction), hitRecord.normal) / math.Sqrt(r.direction.X*r.direction.X+r.direction.Y*r.direction.Y+r.direction.Z*r.direction.Z)
	} else {
		outwardNormal = hitRecord.normal
		niOverNt = 1.0 / d.RefractiveIndex
		cosine = -1 * r3.Dot(r3.Unit(r.direction), hitRecord.normal) / math.Sqrt(r.direction.X*r.direction.X+r.direction.Y*r.direction.Y+r.direction.Z*r.direction.Z)
	}

	shouldRefract, refractedVec := refracted(&r.direction, &outwardNormal, niOverNt)
	if shouldRefract {
		reflectProb = schlick(cosine, d.RefractiveIndex)
	} else {
		reflectProb = 1.0
	}

	if rand.Float64() < reflectProb {
		return true, r3.Vec{X: 1.0, Y: 1.0, Z: 1.0}, ray{p: r3.Add(hitRecord.p, r3.Scale(0.00001, reflectedVec)), direction: reflectedVec}, r3.Vec{}
	} else {
		return true, r3.Vec{X: 1.0, Y: 1.0, Z: 1.0}, ray{p: r3.Add(hitRecord.p, r3.Scale(0.00001, refractedVec)), direction: refractedVec}, r3.Vec{}
	}
}

// see https://www.cs.uregina.ca/Links/class-info/315/WWW/Lab4/#Lighting
func (p PhongBlinn) scatter(is *ImageSpec, r *ray, hitRecord *hitRecord, bvh *boundingVolumeHierarchy, lights *[]Light) (shouldTrace bool, attenuation r3.Vec, scattered ray, color r3.Vec) {
	c := r3.Vec{}
	for _, light := range *lights {
		if light.hasPosition() {
			monteCarloRepetitions := is.SoftShadowMonteCarloRepetitions
			monteCarloMaxLength := softShadowMonteCarloMaxLengthDeviation
			for i := 0; i < monteCarloRepetitions; i++ {
				hitPoint := hitRecord.p
				monteCarloVariance := r3.Scale(monteCarloMaxLength, randomInUnitSphere())
				if light.isPointVisible(&hitPoint, bvh, &monteCarloVariance) {
					lightPosition := *light.getPosition()
					lightToPoint := r3.Sub(lightPosition, hitPoint)
					lightDirection := r3.Unit(lightToPoint)
					lightDistanceSqrd := lightToPoint.X*lightToPoint.X + lightToPoint.Y*lightToPoint.Y + lightToPoint.Z*lightToPoint.Z

					// diffuse Color merges lighting Color and material Color
					nDotL := r3.Dot(hitRecord.normal, lightDirection)
					intensity := saturate(nDotL)
					lightColor := light.getColorFrac()
					diffuseColor := r3.Scale(
						intensity*light.getLightIntensity()/lightDistanceSqrd,
						r3.Unit(r3.Vec{X: p.Color.X + lightColor.X, Y: p.Color.Y + lightColor.Y, Z: p.Color.Z + lightColor.Z}),
					)

					// specular Color uses specular Color of material
					h := r3.Unit(r3.Add(lightDirection, r3.Unit(r.direction)))
					nDotH := r3.Dot(hitRecord.normal, h)
					specIntensity := math.Pow(saturate(nDotH), p.SpecHardness)
					specularColor := r3.Scale(specIntensity*light.getSpecularLightIntensity()/lightDistanceSqrd, p.SpecularColor)

					combinedColor := r3.Vec{
						X: math.Min(1.0, diffuseColor.X+specularColor.X),
						Y: math.Min(1.0, diffuseColor.Y+specularColor.Y),
						Z: math.Min(1.0, diffuseColor.Z+specularColor.Z),
					}
					c = r3.Add(c, r3.Scale(1/float64(monteCarloRepetitions), combinedColor))
				}
			}
		} else {
			// ambient light merges lighting Color and material Color
			c = r3.Add(c, r3.Scale(light.getLightIntensity(), r3.Unit(r3.Add(p.Color, light.getColorFrac()))))
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

// keeps integer between 0-1 inclusive
func saturate(i float64) float64 {
	if i > 1 {
		return 1
	}
	if i < 0 {
		return 0
	}
	return i
}
