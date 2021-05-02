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
	panic("implement me")
}

func (s sphere) centroid() r3.Vec {
	return s.center
}

func (tr triangle) hit(r *ray, tMin float64, tMax float64) hitRecord {
	bMinusA := r3.Sub(tr.pointB, tr.pointA)
	cMinusA := r3.Sub(tr.pointC, tr.pointA)
	normal := r3.Unit(r3.Cross(bMinusA, cMinusA))
	// area := math.Sqrt(normal.X*normal.X + normal.Y*normal.Y + normal.Z*normal.Z)

	// check if ray and plane are parallel
	nDotRayDirection := r3.Dot(normal, r.direction)
	if math.Abs(nDotRayDirection) < 0.00001 {
		return hitRecord{t: -1}
	}
	// check for backward facing triangle
	if tr.singleSided && r3.Dot(r.direction, normal) > 0 {
		return hitRecord{t: -1}
	}

	// compute d parameter in plane equation
	d := r3.Dot(normal, tr.pointA)

	// compute t
	t := (d - r3.Dot(normal, r.p)) / nDotRayDirection
	// check if the triangle is in behind the ray
	if t < tMin || t > tMax {
		return hitRecord{t: -1}
	}

	// compute the intersection point using ray equation
	p := r.PointAtT(t)

	// Step 2: inside-outside test
	var c r3.Vec // vector perpendicular to triangle's plane

	// edge 0
	edge0 := r3.Sub(tr.pointB, tr.pointA)
	vp0 := r3.Sub(p, tr.pointA)
	c = r3.Cross(edge0, vp0)
	if r3.Dot(normal, c) < 0 {
		return hitRecord{t: -1} // p is on the right side
	}

	// edge 1
	edge1 := r3.Sub(tr.pointC, tr.pointB)
	vp1 := r3.Sub(p, tr.pointB)
	c = r3.Cross(edge1, vp1)
	if r3.Dot(normal, c) < 0 {
		return hitRecord{t: -1} // p is on the right side
	}

	// edge 2
	edge2 := r3.Sub(tr.pointA, tr.pointC)
	vp2 := r3.Sub(p, tr.pointC)
	c = r3.Cross(edge2, vp2)
	if r3.Dot(normal, c) < 0 {
		return hitRecord{t: -1} // p is on the right side
	}

	return hitRecord{
		t:        t,
		p:        p,
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
