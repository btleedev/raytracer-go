package raytracer

import (
	"gonum.org/v1/gonum/spatial/r3"
	"math"
)

type hitRecord struct {
	t        float64
	p        r3.Vec
	normal   r3.Vec
	material *material
}

type shape interface {
	hit(r *ray, tMin float64, tMax float64) hitRecord
	translate(tv r3.Vec)
	scale(c float64)
	// rotation vector is in degrees
	rotate(rv r3.Vec)
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
				material: &s.mat,
			}
		}
		secondPoint := (-b - math.Sqrt(b*b-a*c)) / a
		if secondPoint > tMin && firstPoint <= tMax {
			return hitRecord{
				t:        secondPoint,
				p:        r.PointAtT(secondPoint),
				normal:   r3.Scale(1/s.radius, r3.Sub(r.PointAtT(secondPoint), s.center)),
				material: &s.mat,
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
		material: &tr.mat,
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

func (b boundingBox) hit(r *ray, tMin float64, tMax float64) hitRecord {
	normalizedDir := r3.Unit(r.direction)
	invDirection := r3.Vec{
		X: 1 / normalizedDir.X,
		Y: 1 / normalizedDir.Y,
		Z: 1 / normalizedDir.Z,
	}
	// 1 if less than 0, invert if less than 0
	bounds0 := b.pMin
	bounds1 := b.pMax
	if r.direction.X < 0 {
		bounds0.X = b.pMax.X
		bounds1.X = b.pMin.X
	}
	if r.direction.Y < 0 {
		bounds0.Y = b.pMax.Y
		bounds1.Y = b.pMin.Y
	}
	if r.direction.Z < 0 {
		bounds0.Z = b.pMax.Z
		bounds1.Z = b.pMin.Z
	}

	ptMin := (bounds0.X - r.p.X) * invDirection.X
	ptMax := (bounds1.X - r.p.X) * invDirection.X
	tYMin := (bounds0.Y - r.p.Y) * invDirection.Y
	tYMax := (bounds1.Y - r.p.Y) * invDirection.Y

	if ptMin > tYMax || tYMin > ptMax {
		return hitRecord{t: -1}
	}

	if tYMin > ptMin {
		ptMin = tYMin
	}
	if tYMax < ptMax {
		ptMax = tYMax
	}

	tZMin := (bounds0.Z - r.p.Z) * invDirection.Z
	tZMax := (bounds1.Z - r.p.Z) * invDirection.Z

	if (ptMin > tZMax) || (ptMax < tZMin) {
		return hitRecord{t: -1}
	}

	if tZMin > ptMin {
		ptMin = tZMin
	}
	if tZMax < ptMax {
		ptMax = tZMax
	}

	tHit := tMin
	if tHit < tMin || tHit > tMax {
		return hitRecord{t: -1}
	}

	_, hr := trace(r, &(b.shapes), tMin)
	return *hr
}

func (b *boundingBox) translate(tv r3.Vec) {
	for _, v := range b.shapes {
		v.translate(tv)
	}
}

func (b *boundingBox) scale(c float64) {
	for _, v := range b.shapes {
		v.scale(c)
	}
}

func (b *boundingBox) rotate(rv r3.Vec) {
	for _, v := range b.shapes {
		v.rotate(rv)
	}
}
