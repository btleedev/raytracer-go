package raytracer

import (
	"container/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    *boundingVolumeHierarchyNode // The value of the item; arbitrary.
	priority float64                      // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A bvhPriorityQueue implements heap.Interface and holds Items.
type bvhPriorityQueue []*Item

func (pq bvhPriorityQueue) Len() int { return len(pq) }

func (pq bvhPriorityQueue) Less(i, j int) bool {
	// give priority to the lowest
	return pq[i].priority < pq[j].priority
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

// update modifies the priority and value of an Item in the queue.
func (pq *bvhPriorityQueue) update(item *Item, value *boundingVolumeHierarchyNode, priority float64) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
