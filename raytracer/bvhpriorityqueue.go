package raytracer

import (
	"container/heap"
	"math"
)

// An Item is something we manage in a t queue.
type Item struct {
	value *boundingVolumeHierarchyNode // The value of the item; arbitrary.
	// Value of 't' in the Ray equation.
	t float64
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A bvhPriorityQueue implements heap.Interface and holds Items.
type bvhPriorityQueue []*Item

func (pq bvhPriorityQueue) Len() int { return len(pq) }

func (pq bvhPriorityQueue) Less(i, j int) bool {
	// give priority to closest t, if the same give priority to node id order
	if math.Abs(pq[i].t-pq[j].t) < 1e-9 {
		return pq[i].value.nodeId < pq[j].value.nodeId
	}
	return pq[i].t < pq[j].t
}

func (pq bvhPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *bvhPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *bvhPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the t and value of an Item in the queue.
func (pq *bvhPriorityQueue) update(item *Item, value *boundingVolumeHierarchyNode, priority float64) {
	item.value = value
	item.t = priority
	heap.Fix(pq, item.index)
}
