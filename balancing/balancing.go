package balancing

import (
	"math/rand"
	"proj3/deque"
)

type WorkSharingThread struct{
	queue 		[]deque.Deque //each thread has its own deque of tasks kept in an array shared by all threads
	random 		rand.Rand
	threshold 	int
}

func NewWorkSharingThread(deqSlice []deque.Deque, threshold int) *WorkSharingThread{
	return &WorkSharingThread{deqSlice, rand.Rand{}, threshold}
}

func (ws WorkSharingThread) Run(goID int){
	me := goID
	size := ws.queue[me].Length()
	if rand.Intn(size + 1) == size {
		victim := rand.Intn(len(ws.queue))
		min := 0
		max := 0
		if victim <= me {
			min = victim
			max = me
		}else{
			min = me
			max = victim
		}
		if me == victim{
			return
		}
		ws.queue[min].Lock()
		ws.queue[max].Lock()
		ws.balance(&ws.queue[min], &ws.queue[max])
		ws.queue[max].Unlock()
		ws.queue[min].Unlock()
	}
}

func (ws WorkSharingThread) balance(d0 *deque.Deque, d1 *deque.Deque){
	qMin := &deque.Deque{}
	qMax := &deque.Deque{}

	if d0.Length() < d1.Length(){
		qMin = d0
		qMax = d1
	}else{
		qMin = d1
		qMax = d0
	}

	diff := qMax.Length() - qMin.Length()

	if diff > ws.threshold{
		for qMax.Length() > qMin.Length(){
			success, request := qMax.PopBack()
			if success{
				qMin.PushBack(request)
			}
		}
	}
}