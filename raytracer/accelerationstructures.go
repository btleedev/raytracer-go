package raytracer

import (
	"fmt"
	"gonum.org/v1/gonum/spatial/r3"
	"math"
	"math/rand"
)

type boundingVolumeHierarchyNode struct {
	pMin     r3.Vec
	pMax     r3.Vec
	leaf     bool
	shape    *shape
	children []*boundingVolumeHierarchyNode
}

type boundingVolumeHierarchy struct {
	root   boundingVolumeHierarchyNode
	shapes *[]shape
}

// bounding box hierarchy where boundaries are computed in a box shape
func NewBoundingVolumeHierarchy(shapes *[]shape) *boundingVolumeHierarchy {
	fmt.Printf("Building BoundingVolumeHierarchy\n")
	pMin := r3.Vec{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	pMax := r3.Vec{X: float64(math.MinInt64), Y: float64(math.MinInt64), Z: float64(math.MinInt64)}
	for _, s := range *shapes {
		lowest, highest := s.computeSquareBounds()
		pMin.X = math.Min(pMin.X, lowest.X)
		pMin.Y = math.Min(pMin.Y, lowest.Y)
		pMin.Z = math.Min(pMin.Z, lowest.Z)
		pMax.X = math.Max(pMax.X, highest.X)
		pMax.Y = math.Max(pMax.Y, highest.Y)
		pMax.Z = math.Max(pMax.Z, highest.Z)
	}
	// add the max jitter than can happen when jittering the centroid of shapes
	pMin = r3.Sub(pMin, r3.Scale(bvhCentroidJitterFactor, r3.Vec{X: 1, Y: 1, Z: 1}))
	pMax = r3.Add(pMax, r3.Scale(bvhCentroidJitterFactor, r3.Vec{X: 1, Y: 1, Z: 1}))

	bvh := boundingVolumeHierarchy{
		shapes: shapes,
		root: boundingVolumeHierarchyNode{
			pMin:     pMin,
			pMax:     pMax,
			leaf:     true,
			shape:    nil,
			children: nil,
		},
	}

	for i := 0; i < len(*shapes); i++ {
		ptr := &(*shapes)[i]
		addToBVH(&bvh.root, ptr)
	}
	bvh.recomputeBounds()

	fmt.Printf("Finished building BoundingVolumeHierarchy\n")
	return &bvh
}

func (bvh boundingVolumeHierarchy) trace(r *ray, tMin float64) (hit bool, record *hitRecord) {
	return traceDownBoundingVolumeHierarchyNode(r, tMin, math.MaxFloat64, &bvh.root)
}

// traces a ray and returns if it hits something, and a hit record
func traceDownBoundingVolumeHierarchyNode(r *ray, tMin float64, tMax float64, node *boundingVolumeHierarchyNode) (hit bool, record *hitRecord) {
	if !hitBoundingBox(r, node.pMin, node.pMax, tMin, tMax) {
		return false, &hitRecord{t: -1}
	}

	if node.leaf {
		if node.shape == nil {
			return false, nil
		} else {
			hr := (*node.shape).hit(r, tMin, tMax)
			return hr.t > 0.0, &hr
		}
	} else {
		localTMax := tMax
		var minHitRecord = &hitRecord{
			t: tMax,
		}
		if node.children != nil {
			for _, v := range node.children {
				rHit, rhr := traceDownBoundingVolumeHierarchyNode(r, tMin, localTMax, v)
				if rHit {
					if rhr.t > tMin && rhr.t < minHitRecord.t {
						minHitRecord = rhr
					}
				}
			}
		}
		return minHitRecord.t != math.MaxFloat64, minHitRecord
	}
}

// recomputes the bounds for all objects in the BVH, from bottom up
func (bvh boundingVolumeHierarchy) recomputeBounds() {
	recomputeNodeBounds(&bvh.root)
}

func recomputeNodeBounds(node *boundingVolumeHierarchyNode) (pMin r3.Vec, pMax r3.Vec) {
	boundsLow := r3.Vec{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	boundsHigh := r3.Vec{X: float64(math.MinInt64), Y: float64(math.MinInt64), Z: float64(math.MinInt64)}
	if node.leaf {
		if node.shape != nil {
			boundsLow, boundsHigh = (*node.shape).computeSquareBounds()
		}
	} else {
		for _, child := range node.children {
			childBoundsLow, childBoundsHigh := recomputeNodeBounds(child)
			boundsLow.X = math.Min(boundsLow.X, childBoundsLow.X)
			boundsLow.Y = math.Min(boundsLow.Y, childBoundsLow.Y)
			boundsLow.Z = math.Min(boundsLow.Z, childBoundsLow.Z)
			boundsHigh.X = math.Max(boundsHigh.X, childBoundsHigh.X)
			boundsHigh.Y = math.Max(boundsHigh.Y, childBoundsHigh.Y)
			boundsHigh.Z = math.Max(boundsHigh.Z, childBoundsHigh.Z)
		}
	}

	node.pMin = boundsLow
	node.pMax = boundsHigh
	return node.pMin, node.pMax
}

func addToBVH(
	curr *boundingVolumeHierarchyNode,
	shape *shape,
) {
	if curr.leaf {
		// empty leaf node, feel free to add
		if curr.shape == nil {
			ptr := &(*shape)
			curr.shape = ptr
			return
			// promote this to a child node, put object 1 is 1 and another in 2
		} else {
			curr.leaf = false
			curr.children = splitBvhQuadrant(&curr.pMin, &curr.pMax)
			removedShape := *curr.shape
			curr.shape = nil

			// recursive call to same node, now that it isn't a leaf it should add it
			addToBVH(curr, &removedShape)
			addToBVH(curr, shape)
			return
		}
	} else {
		// delegate adding it to the node down
		ptr := curr.children[getBvhQuadrantIndex(shape, &curr.pMin, &curr.pMax)]
		addToBVH(ptr, shape)
		return
	}
}

// front bottom left = 0
// front bottom right = 1
// front top left = 2
// front top right = 3
// back bottom left = 4
// back bottom right = 5
// back top left = 6
// back top right = 7
// to prevent two shapes from having the same centroid coordinates, we add a random jitter factor to each centroid
func getBvhQuadrantIndex(s *shape, pMin *r3.Vec, pMax *r3.Vec) uint8 {
	centroid := r3.Add((*s).centroid(), r3.Scale(bvhCentroidJitterFactor, r3.Vec{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()}))
	idx := uint8(0)
	if centroid.X > pMin.X+(pMax.X-pMin.X)/2 {
		idx += 1
	}
	if centroid.Y > pMin.Y+(pMax.Y-pMin.Y)/2 {
		idx += 2
	}
	if centroid.Z > pMin.Z+(pMax.Z-pMin.Z)/2 {
		idx += 4
	}
	return idx
}

// see getBvhQuadrantIndex
func splitBvhQuadrant(lowestBounds *r3.Vec, highestBounds *r3.Vec) []*boundingVolumeHierarchyNode {
	halfX := (highestBounds.X - lowestBounds.X) / 2
	halfY := (highestBounds.Y - lowestBounds.Y) / 2
	halfZ := (highestBounds.Z - lowestBounds.Z) / 2
	return []*boundingVolumeHierarchyNode{
		&boundingVolumeHierarchyNode{
			pMin:     r3.Vec{X: lowestBounds.X, Y: lowestBounds.Y, Z: lowestBounds.Z},
			pMax:     r3.Vec{X: lowestBounds.X + halfX, Y: lowestBounds.Y + halfY, Z: lowestBounds.Z + halfZ},
			leaf:     true,
			shape:    nil,
			children: nil,
		},
		&boundingVolumeHierarchyNode{
			pMin:     r3.Vec{X: lowestBounds.X + halfX, Y: lowestBounds.Y, Z: lowestBounds.Z},
			pMax:     r3.Vec{X: highestBounds.X, Y: lowestBounds.Y + halfY, Z: lowestBounds.Z + halfZ},
			leaf:     true,
			shape:    nil,
			children: nil,
		},
		&boundingVolumeHierarchyNode{
			pMin:     r3.Vec{X: lowestBounds.X, Y: lowestBounds.Y + halfY, Z: lowestBounds.Z},
			pMax:     r3.Vec{X: lowestBounds.X + halfX, Y: highestBounds.Y, Z: lowestBounds.Z + halfZ},
			leaf:     true,
			shape:    nil,
			children: nil,
		},
		&boundingVolumeHierarchyNode{
			pMin:     r3.Vec{X: lowestBounds.X + halfX, Y: lowestBounds.Y + halfY, Z: lowestBounds.Z},
			pMax:     r3.Vec{X: highestBounds.X, Y: highestBounds.Y, Z: lowestBounds.Z + halfZ},
			leaf:     true,
			shape:    nil,
			children: nil,
		},
		&boundingVolumeHierarchyNode{
			pMin:     r3.Vec{X: lowestBounds.X, Y: lowestBounds.Y, Z: lowestBounds.Z + halfZ},
			pMax:     r3.Vec{X: lowestBounds.X + halfX, Y: lowestBounds.Y + halfY, Z: highestBounds.Z},
			leaf:     true,
			shape:    nil,
			children: nil,
		},
		&boundingVolumeHierarchyNode{
			pMin:     r3.Vec{X: lowestBounds.X + halfX, Y: lowestBounds.Y, Z: lowestBounds.Z + halfZ},
			pMax:     r3.Vec{X: highestBounds.X, Y: lowestBounds.Y + halfY, Z: highestBounds.Z},
			leaf:     true,
			shape:    nil,
			children: nil,
		},
		&boundingVolumeHierarchyNode{
			pMin:     r3.Vec{X: lowestBounds.X, Y: lowestBounds.Y + halfY, Z: lowestBounds.Z + halfZ},
			pMax:     r3.Vec{X: lowestBounds.X + halfX, Y: highestBounds.Y, Z: highestBounds.Z},
			leaf:     true,
			shape:    nil,
			children: nil,
		},
		&boundingVolumeHierarchyNode{
			pMin:     r3.Vec{X: lowestBounds.X + halfX, Y: lowestBounds.Y + halfY, Z: lowestBounds.Z + halfZ},
			pMax:     r3.Vec{X: highestBounds.X, Y: highestBounds.Y, Z: highestBounds.Z},
			leaf:     true,
			shape:    nil,
			children: nil,
		},
	}
}

// determines whether the ray hits the bounding box
func hitBoundingBox(r *ray, pMin r3.Vec, pMax r3.Vec, tMin float64, tMax float64) bool {
	normalizedDir := r3.Unit(r.direction)
	invDirection := r3.Vec{
		X: 1 / normalizedDir.X,
		Y: 1 / normalizedDir.Y,
		Z: 1 / normalizedDir.Z,
	}
	// 1 if less than 0, invert if less than 0
	bounds0 := pMin
	bounds1 := pMax
	if r.direction.X < 0 {
		bounds0.X = pMax.X
		bounds1.X = pMin.X
	}
	if r.direction.Y < 0 {
		bounds0.Y = pMax.Y
		bounds1.Y = pMin.Y
	}
	if r.direction.Z < 0 {
		bounds0.Z = pMax.Z
		bounds1.Z = pMin.Z
	}

	ptMin := (bounds0.X - r.p.X) * invDirection.X
	ptMax := (bounds1.X - r.p.X) * invDirection.X
	tYMin := (bounds0.Y - r.p.Y) * invDirection.Y
	tYMax := (bounds1.Y - r.p.Y) * invDirection.Y

	if ptMin > tYMax || tYMin > ptMax {
		return false
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
		return false
	}

	if tZMin > ptMin {
		ptMin = tZMin
	}
	if tZMax < ptMax {
		ptMax = tZMax
	}

	tHit := tMin
	if tHit < tMin || tHit > tMax {
		return false
	}

	return true
}
