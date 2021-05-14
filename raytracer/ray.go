package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
)

type ray struct {
	p                   r3.Vec
	normalizedDirection r3.Vec
}

func (r ray) PointAtT(t float64) r3.Vec {
	return r3.Add(r.p, r3.Scale(t, r.normalizedDirection))
}
