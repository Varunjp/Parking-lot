package usecase

import (
	"container/heap"
	"parking-lot/internal/domain"
	"sync"
)

type Dispatcher struct {
	queue     *PQ
	allocator *Allocator
	mu        sync.Mutex
}

type Result struct {
	Level  int 
	Slot   int 
	Err    error 
}

func NewDispatcher(a *Allocator) *Dispatcher {
	pq := &PQ{}
	pq.Init()
	return &Dispatcher{queue: pq, allocator: a}
}

func (d *Dispatcher) AddRequest(v domain.Vehicle) Result{
	priority,err := getPriority(v.CustomerType)
	if err != nil {
		return Result{
			Level: 0,
			Slot: 0,
			Err: err,
		} 
	}
	respChan := make(chan Result)
	req := Request {
		Vehicle: v,
		Priority: priority,
		RespChan: respChan,
	}
	d.mu.Lock()
	heap.Push(d.queue,req)
	d.mu.Unlock()
	return <-respChan
}

func (d *Dispatcher) Process() {
	for d.queue.Len() > 0 {
		d.mu.Lock()
		req := heap.Pop(d.queue).(Request)
		d.mu.Unlock()
		v := req.Vehicle.(domain.Vehicle)

		level,slot,err := d.allocator.Allocate(v)
		
		req.RespChan <- Result{
			Level: level,
			Slot: slot,
			Err: err,
		}
	}
}

func getPriority(t domain.CustomerType) (int,error) {
	switch t {
	case domain.Emergency:
		return 1,nil
	case domain.VIP:
		return 2,nil
	case domain.Regular:
		return 3,nil
	}
	return 0,domain.ErrInvalidPriority
}