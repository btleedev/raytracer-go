package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
	"math/rand"
)

type camera struct{
	w, u, v r3.Vec
	origin r3.Vec
	lowerLeftCorner r3.Vec
	horizontal r3.Vec
	vertical r3.Vec
	lensRadius float64
}

func NewCamera(lookFrom r3.Vec, lookAt r3.Vec, up r3.Vec, fov float64, aspect float64, aperature float64, focusDist float64) camera {
	theta := fov * math.Pi / 180.0
	halfHeight := math.Tan(theta/2)
	halfWidth := aspect * halfHeight
	w := r3.Unit(r3.Sub(lookFrom, lookAt))
	u := r3.Unit(r3.Cross(up, w))
	v := r3.Cross(w, u)
	return camera{
		w: 					w,
		u:					u,
		v:					v,
		origin:          	lookFrom,
		lowerLeftCorner: 	r3.Sub(r3.Sub(r3.Sub(lookFrom, r3.Scale(halfWidth * focusDist, u)), r3.Scale(halfHeight * focusDist, v)), r3.Scale(focusDist, w)),
		horizontal:      	r3.Scale(2 * halfWidth * focusDist, u),
		vertical:      		r3.Scale(2 * halfHeight * focusDist, v),
		lensRadius: 		aperature / 2,
	}
}

func (c camera) getRay(s float64, t float64) ray {
	rd := r3.Scale(c.lensRadius, randomInUnitDisk())
	offset := r3.Add(r3.Scale(rd.X, c.u), r3.Scale(rd.Y, c.v))
	return ray{
		p:       	r3.Add(c.origin, offset),
		direction: 	r3.Sub(r3.Sub(r3.Add(r3.Add(c.lowerLeftCorner, r3.Scale(s, c.horizontal)), r3.Scale(t, c.vertical)), c.origin), offset),
	}
}

func randomInUnitDisk() r3.Vec {
	p := r3.Vec{}
	for {
		p = r3.Sub(r3.Scale(2, r3.Vec{ X: rand.Float64(), Y: rand.Float64(), Z: 0 }), r3.Vec{ X: 1, Y: 1, Z: 0 })
		if r3.Dot(p, p) < 1.0 {
			break
		}
	}
	return p
}