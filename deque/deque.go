package deque

import (
	"fmt"
	"sync"
)

//Implementation of a double ended queue

type Deque struct{
	Items []interface{}
	mu 		sync.Mutex
}

// Create new deque
func NewDeque() *Deque{
	i := make([]interface{}, 0)
	return &Deque{i, sync.Mutex{}}
}

// Checks whether Deque is empty
func (d *Deque) IsEmpty() bool {
	return len(d.Items) == 0
}

// Get length of deque
func (d *Deque) Length() int {
	return len(d.Items)
}

// Worker adds/ removes from the Front

// Remove an item from the Front
func (d *Deque) PopFront() (bool, interface{}){
	if len(d.Items) > 0 {
		temp := d.Items[0]
		d.Items = d.Items[1:]
		return true, temp
	}
	return false, nil
}

// Add an item to the Front
func (d *Deque) PushFront(item interface{}) {
	temp := []interface{}{item}
	d.Items = append(temp, d.Items ...)
}

// Thief adds/ removes from the Back

// Remove an item from the Back
func (d *Deque) PopBack() (bool, interface{}){
	if len(d.Items) > 0 {
		l := len(d.Items) - 1
		temp := d.Items[l]
		d.Items = d.Items[:l-1]
		return true, temp
	}
	return false, nil
}

// Add an item to the Back
func (d *Deque) PushBack(item interface{}) {
	d.Items = append(d.Items, item)
}

// Locks access
func (d *Deque) Lock() {
	d.mu.Lock()
}

// Unlocks access
func (d *Deque) Unlock() {
	d.mu.Unlock()
}

// Print deque
func (d *Deque) Print() {
	fmt.Println(d)
}


