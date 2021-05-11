package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
)

type hitRecord struct {
	t        float64
	p        r3.Vec
	normal   r3.Vec
	material Material
}

type Shape interface {
	hit(r *ray, tMin float64, tMax float64) hitRecord
	translate(tv r3.Vec)
	scale(c float64)
	// rotation vector is in degrees
	rotate(rv r3.Vec)
	computeSquareBounds() (lowest r3.Vec, highest r3.Vec)
	centroid() r3.Vec
}

type Sphere struct {
	Center r3.Vec
	Radius float64
	Mat    Material
}

type Triangle struct {
	PointA      r3.Vec
	PointB      r3.Vec
	PointC      r3.Vec
	SingleSided bool
	Mat         Material
}

func (s Sphere) hit(r *ray, tMin float64, tMax float64) hitRecord {
	oc := r3.Sub(r.p, s.Center)
	a := r3.Dot(r.direction, r.direction)
	b := r3.Dot(oc, r.direction)
	c := r3.Dot(oc, oc) - s.Radius*s.Radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		firstPoint := (-b - math.Sqrt(b*b-a*c)) / a
		if firstPoint > tMin && firstPoint <= tMax {
			return hitRecord{
				t:        firstPoint,
				p:        r.PointAtT(firstPoint),
				normal:   r3.Scale(1/s.Radius, r3.Sub(r.PointAtT(firstPoint), s.Center)),
				material: s.Mat,
			}
		}
		secondPoint := (-b - math.Sqrt(b*b-a*c)) / a
		if secondPoint > tMin && firstPoint <= tMax {
			return hitRecord{
				t:        secondPoint,
				p:        r.PointAtT(secondPoint),
				normal:   r3.Scale(1/s.Radius, r3.Sub(r.PointAtT(secondPoint), s.Center)),
				material: s.Mat,
			}
		}
	}
	return hitRecord{
		t: -1,
	}
}

func (s *Sphere) translate(tv r3.Vec) {
	s.Center = r3.Add(tv, s.Center)
}

func (s *Sphere) scale(c float64) {
	s.Radius *= c
}

func (s *Sphere) rotate(rv r3.Vec) {
}

func (s Sphere) computeSquareBounds() (lowest r3.Vec, highest r3.Vec) {
	return r3.Sub(s.Center, r3.Vec{X: s.Radius, Y: s.Radius, Z: s.Radius}), r3.Add(s.Center, r3.Vec{X: s.Radius, Y: s.Radius, Z: s.Radius})
}

func (s Sphere) centroid() r3.Vec {
	return s.Center
}

func (tr Triangle) hit(r *ray, tMin float64, tMax float64) hitRecord {
	// moller-trumbore ray triangle intersection algorithm
	dir := r.direction
	bMinusA := r3.Sub(tr.PointB, tr.PointA)
	cMinusA := r3.Sub(tr.PointC, tr.PointA)
	normal := r3.Unit(r3.Cross(bMinusA, cMinusA))
	pvec := r3.Cross(dir, cMinusA)
	det := r3.Dot(bMinusA, pvec)

	if tr.SingleSided {
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

	tvec := r3.Sub(r.p, tr.PointA)
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
		material: tr.Mat,
	}
}

func (t *Triangle) translate(tv r3.Vec) {
	t.PointA = r3.Add(tv, t.PointA)
	t.PointB = r3.Add(tv, t.PointB)
	t.PointC = r3.Add(tv, t.PointC)
}

func (t *Triangle) scale(c float64) {
	t.PointA = r3.Scale(c, t.PointA)
	t.PointB = r3.Scale(c, t.PointB)
	t.PointC = r3.Scale(c, t.PointC)
}

func (t *Triangle) rotate(rv r3.Vec) {
	t.PointA = rotatePoint(t.PointA, rv)
	t.PointB = rotatePoint(t.PointB, rv)
	t.PointC = rotatePoint(t.PointC, rv)
}

func (tr Triangle) computeSquareBounds() (lowest r3.Vec, highest r3.Vec) {
	pMin := r3.Vec{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	pMax := r3.Vec{X: float64(math.MinInt64), Y: float64(math.MinInt64), Z: float64(math.MinInt64)}

	pMin.X = math.Min(pMin.X, tr.PointA.X)
	pMin.X = math.Min(pMin.X, tr.PointB.X)
	pMin.X = math.Min(pMin.X, tr.PointC.X)
	pMin.Y = math.Min(pMin.Y, tr.PointA.Y)
	pMin.Y = math.Min(pMin.Y, tr.PointB.Y)
	pMin.Y = math.Min(pMin.Y, tr.PointC.Y)
	pMin.Z = math.Min(pMin.Z, tr.PointA.Z)
	pMin.Z = math.Min(pMin.Z, tr.PointB.Z)
	pMin.Z = math.Min(pMin.Z, tr.PointC.Z)

	pMax.X = math.Max(pMax.X, tr.PointA.X)
	pMax.X = math.Max(pMax.X, tr.PointB.X)
	pMax.X = math.Max(pMax.X, tr.PointC.X)
	pMax.Y = math.Max(pMax.Y, tr.PointA.Y)
	pMax.Y = math.Max(pMax.Y, tr.PointB.Y)
	pMax.Y = math.Max(pMax.Y, tr.PointC.Y)
	pMax.Z = math.Max(pMax.Z, tr.PointA.Z)
	pMax.Z = math.Max(pMax.Z, tr.PointB.Z)
	pMax.Z = math.Max(pMax.Z, tr.PointC.Z)
	return pMin, pMax
}

func (tr Triangle) centroid() r3.Vec {
	return r3.Scale(1/3.0, r3.Add(tr.PointA, r3.Add(tr.PointB, tr.PointC)))
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
