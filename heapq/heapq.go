package heapq

import (
	"container/heap"
)

// Contains the price to get to this city from the source
// and the path corresponding to that price
type Path struct{
	City string 		//Current city
	Price int			//Cheapest price from source to this city
	Cities []string 	//Path order from source to this city
}

// Priority queue ordering cities by price
type PriorityQueue []Path

// Returns length of the priority queue
func (pq PriorityQueue) Len() int {
	return len(pq)
}

// Method needed to classify PriorityQueue as heap.Interface
func (pq PriorityQueue) Less(i, j int) bool {
	//Min heap - want cheapest flight first
	return pq[i].Price < pq[j].Price
}

// Method needed to classify PriorityQueue as heap.Interface
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push new path into priority queue
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Path))
}

// Pop cheapest path out from priority queue
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

//Heap that push/pops from the priority queue
type Heapq struct {
	PQ	*PriorityQueue
}

// Returns pointer to priority queue
func NewHeapQ() *Heapq{
	return &Heapq{&PriorityQueue{}}
}

// Push new path into priority queue
func (hq *Heapq) Push(p Path){
	heap.Push(hq.PQ, p)
}

//Return the path with the cheapest price
func (hq *Heapq) Pop() Path{
	for hq.PQ.Len() > 0 {
		x := heap.Pop(hq.PQ).(Path)
		return x
	}
	return Path{}
}


