package queue

import (
	"fmt"
	"sync"
)

// Implementation of a coarse grained queue

type Queue struct{
	Items 	[]interface{}
	lock 	sync.Mutex
}

// Create new deque
func NewQueue() *Queue{
	return &Queue{nil, sync.Mutex{}}
}

// Checks whether Deque is empty
func (q *Queue) IsEmpty() bool {
	return len(q.Items) == 0
}

// Get length of deque
func (q *Queue) Length() int {
	return len(q.Items)
}


// Remove an item from the Front
func (q *Queue) Dequeue() interface{}{

	if len(q.Items) > 0 {
		temp := q.Items[0]
		q.Items = q.Items[1:]
		return temp
	}
	return nil
}

// Add an item to the Back
func (q *Queue) Enqueue(item interface{}) {
	q.Items = append(q.Items, item)
}

func (q *Queue) Display() {
	fmt.Println(q)
}

func (q *Queue) Lock() {
	q.lock.Lock()
}

func (q *Queue) Unlock() {
	q.lock.Unlock()
}

