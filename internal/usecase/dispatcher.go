package usecase

import (
	"container/heap"
	"parking-lot/internal/domain"
)

type Dispatcher struct {
	queue     *PQ
	allocator *Allocator
}

func NewDispatcher(a *Allocator) *Dispatcher {
	pq := &PQ{}
	pq.Init()
	return &Dispatcher{queue: pq, allocator: a}
}

func (d *Dispatcher) AddRequest(v domain.Vehicle) {
	priority := getPriority(v.CustomerType)
	heap.Push(d.queue,Request{Vehicle: v,Priority: priority})
}

func (d *Dispatcher) Process() {
	for d.queue.Len() > 0 {
		req := heap.Pop(d.queue).(Request)
		v := req.Vehicle.(domain.Vehicle)

		level,slot,err := d.allocator.Allocate(v)
		if err != nil {
			println("Error:",err.Error())
			continue
		}

		println("Allocated:",v.ID," Level:",level," Slot:",slot)
	}
}

func getPriority(t domain.CustomerType) int {
	switch t {
	case domain.Emergency:
		return 1
	case domain.VIP:
		return 2
	default:
		return 3
	}
}