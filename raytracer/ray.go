package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
)

type ray struct {
	p         r3.Vec
	direction r3.Vec
}

func (r ray) PointAtT(t float64) r3.Vec {
	return r3.Add(r.p, r3.Scale(t, r.direction))
}

func trace(r *ray, shapes *[]shape, tMin float64) (hit bool, record *hitRecord) {
	var minHitRecord = hitRecord{
		t: math.MaxFloat64,
	}
	for _, shape := range *shapes {
		hitRecord := shape.hit(r, tMin, minHitRecord.t)
		if hitRecord.t > 0.0 && hitRecord.t < minHitRecord.t {
			minHitRecord = hitRecord
		}
	}
	return minHitRecord.t != math.MaxFloat64, &minHitRecord
}
