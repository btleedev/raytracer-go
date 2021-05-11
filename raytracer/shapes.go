package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
)

type hitRecord struct {
	t        float64
	p        r3.Vec
	normal   r3.Vec
	material material
}

type shape interface {
	hit(r *ray, tMin float64, tMax float64) hitRecord
	translate(tv r3.Vec)
	scale(c float64)
	// rotation vector is in degrees
	rotate(rv r3.Vec)
	computeSquareBounds() (lowest r3.Vec, highest r3.Vec)
	centroid() r3.Vec
}

type sphere struct {
	center r3.Vec
	radius float64
	mat    material
}

type triangle struct {
	pointA      r3.Vec
	pointB      r3.Vec
	pointC      r3.Vec
	singleSided bool
	mat         material
}

type boundingBox struct {
	pMin   r3.Vec
	pMax   r3.Vec
	shapes []shape
}

func (s sphere) hit(r *ray, tMin float64, tMax float64) hitRecord {
	oc := r3.Sub(r.p, s.center)
	a := r3.Dot(r.direction, r.direction)
	b := r3.Dot(oc, r.direction)
	c := r3.Dot(oc, oc) - s.radius*s.radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		firstPoint := (-b - math.Sqrt(b*b-a*c)) / a
		if firstPoint > tMin && firstPoint <= tMax {
			return hitRecord{
				t:        firstPoint,
				p:        r.PointAtT(firstPoint),
				normal:   r3.Scale(1/s.radius, r3.Sub(r.PointAtT(firstPoint), s.center)),
				material: s.mat,
			}
		}
		secondPoint := (-b - math.Sqrt(b*b-a*c)) / a
		if secondPoint > tMin && firstPoint <= tMax {
			return hitRecord{
				t:        secondPoint,
				p:        r.PointAtT(secondPoint),
				normal:   r3.Scale(1/s.radius, r3.Sub(r.PointAtT(secondPoint), s.center)),
				material: s.mat,
			}
		}
	}
	return hitRecord{
		t: -1,
	}
}

func (s *sphere) translate(tv r3.Vec) {
	s.center = r3.Add(tv, s.center)
}

func (s *sphere) scale(c float64) {
	s.radius *= c
}

func (s *sphere) rotate(rv r3.Vec) {
}

func (s sphere) computeSquareBounds() (lowest r3.Vec, highest r3.Vec) {
	return r3.Sub(s.center, r3.Vec{X: s.radius, Y: s.radius, Z: s.radius}), r3.Add(s.center, r3.Vec{X: s.radius, Y: s.radius, Z: s.radius})
}

func (s sphere) centroid() r3.Vec {
	return s.center
}

func (tr triangle) hit(r *ray, tMin float64, tMax float64) hitRecord {
	// moller-trumbore ray triangle intersection algorithm
	dir := r.direction
	bMinusA := r3.Sub(tr.pointB, tr.pointA)
	cMinusA := r3.Sub(tr.pointC, tr.pointA)
	normal := r3.Unit(r3.Cross(bMinusA, cMinusA))
	pvec := r3.Cross(dir, cMinusA)
	det := r3.Dot(bMinusA, pvec)

	if tr.singleSided {
		if det < 0.0 {
			return hitRecord{t: -1}
		}
	} else {
		// check for parallelism
		if math.Abs(det) < 0.0 {
			return hitRecord{t: -1}
		}
	}

	invDet := 1 / det

	tvec := r3.Sub(r.p, tr.pointA)
	u := r3.Dot(tvec, pvec) * invDet
	if u < 0 || u > 1 {
		return hitRecord{t: -1}
	}

	qvec := r3.Cross(tvec, bMinusA)
	v := r3.Dot(dir, qvec) * invDet
	if v < 0 || u+v > 1 {
		return hitRecord{t: -1}
	}

	t := r3.Dot(cMinusA, qvec) * invDet
	if t < tMin || t > tMax {
		return hitRecord{t: -1}
	}

	return hitRecord{
		t:        t,
		p:        r.PointAtT(t),
		normal:   normal,
		material: tr.mat,
	}
}

func (t *triangle) translate(tv r3.Vec) {
	t.pointA = r3.Add(tv, t.pointA)
	t.pointB = r3.Add(tv, t.pointB)
	t.pointC = r3.Add(tv, t.pointC)
}

func (t *triangle) scale(c float64) {
	t.pointA = r3.Scale(c, t.pointA)
	t.pointB = r3.Scale(c, t.pointB)
	t.pointC = r3.Scale(c, t.pointC)
}

func (t *triangle) rotate(rv r3.Vec) {
	t.pointA = rotatePoint(t.pointA, rv)
	t.pointB = rotatePoint(t.pointB, rv)
	t.pointC = rotatePoint(t.pointC, rv)
}

func (tr triangle) computeSquareBounds() (lowest r3.Vec, highest r3.Vec) {
	pMin := r3.Vec{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	pMax := r3.Vec{X: float64(math.MinInt64), Y: float64(math.MinInt64), Z: float64(math.MinInt64)}

	pMin.X = math.Min(pMin.X, tr.pointA.X)
	pMin.X = math.Min(pMin.X, tr.pointB.X)
	pMin.X = math.Min(pMin.X, tr.pointC.X)
	pMin.Y = math.Min(pMin.Y, tr.pointA.Y)
	pMin.Y = math.Min(pMin.Y, tr.pointB.Y)
	pMin.Y = math.Min(pMin.Y, tr.pointC.Y)
	pMin.Z = math.Min(pMin.Z, tr.pointA.Z)
	pMin.Z = math.Min(pMin.Z, tr.pointB.Z)
	pMin.Z = math.Min(pMin.Z, tr.pointC.Z)

	pMax.X = math.Max(pMax.X, tr.pointA.X)
	pMax.X = math.Max(pMax.X, tr.pointB.X)
	pMax.X = math.Max(pMax.X, tr.pointC.X)
	pMax.Y = math.Max(pMax.Y, tr.pointA.Y)
	pMax.Y = math.Max(pMax.Y, tr.pointB.Y)
	pMax.Y = math.Max(pMax.Y, tr.pointC.Y)
	pMax.Z = math.Max(pMax.Z, tr.pointA.Z)
	pMax.Z = math.Max(pMax.Z, tr.pointB.Z)
	pMax.Z = math.Max(pMax.Z, tr.pointC.Z)
	return pMin, pMax
}

func (tr triangle) centroid() r3.Vec {
	return r3.Scale(1/3.0, r3.Add(tr.pointA, r3.Add(tr.pointB, tr.pointC)))
}

func rotatePoint(point r3.Vec, rv r3.Vec) r3.Vec {
	piDivide180 := math.Pi / 180.0
	rotatedPoint := point

	// around z axis
	x := rotatedPoint.X*math.Cos(piDivide180*rv.Z) - rotatedPoint.Y*math.Sin(piDivide180*rv.Z)
	y := rotatedPoint.X*math.Sin(piDivide180*rv.Z) + rotatedPoint.Y*math.Cos(piDivide180*rv.Z)
	rotatedPoint.X = x
	rotatedPoint.Y = y

	// around x axis
	y = rotatedPoint.Y*math.Cos(piDivide180*rv.X) - rotatedPoint.Z*math.Sin(piDivide180*rv.X)
	z := rotatedPoint.Y*math.Sin(piDivide180*rv.X) + rotatedPoint.Z*math.Cos(piDivide180*rv.X)
	rotatedPoint.Y = y
	rotatedPoint.Z = z

	// around y axis
	x = rotatedPoint.X*math.Cos(piDivide180*rv.Y) + rotatedPoint.Z*math.Sin(piDivide180*rv.Y)
	z = -1*rotatedPoint.X*math.Sin(piDivide180*rv.Y) + rotatedPoint.Z*math.Cos(piDivide180*rv.Y)
	rotatedPoint.X = x
	rotatedPoint.Z = z

	return rotatedPoint
}
