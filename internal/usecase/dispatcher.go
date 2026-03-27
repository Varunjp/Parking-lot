package usecase

import (
	"container/heap"
	"parking-lot/internal/domain"
	"sync"
)

// Dispatcher is responsible for handling incoming parking requests
// and processing them based on priority using a priority queue.
//
// It runs a background worker (Run) that continuously dequeues requests
// and delegates slot allocation to the Allocator.
type Dispatcher struct {
	queue     *PQ  			// Priority queue holding incoming requests
	allocator *Allocator	// Handles actual slot allocation logic
	mu        sync.Mutex	// Protects concurrent access to the queue
}

// Result represents the outcome of a parking allocation request.
type Result struct {
	Level  int 		// Allocated parking level
	Slot   int 		// Allocated parking slot
	Err    error 	// Error if allocation fails
}

// NewDispatcher initializes a Dispatcher and starts the background worker.
//
// NOTE:
// - The Run() method is executed as a goroutine.
// - Dispatcher is designed to be long-lived.
func NewDispatcher(a *Allocator) *Dispatcher {
	pq := &PQ{}
	pq.Init()

	d := &Dispatcher{queue: pq, allocator: a}

	// Start background worker to process requests
	go d.Run()

	return d
}

// AddRequest adds a new vehicle parking request to the queue.
//
// Flow:
// 1. Determine priority based on customer type
// 2. Push request into priority queue
// 3. Wait synchronously for allocation result
//
// This method blocks until the request is processed.
func (d *Dispatcher) AddRequest(v domain.Vehicle) Result{
	
	// Determine priority of the request
	priority,err := getPriority(v.CustomerType)
	if err != nil {
		return Result{Err: err} 
	}
	
	// Channel to receive async response from dispatcher worker
	respChan := make(chan Result)
	
	req := Request {
		Vehicle: v,
		Priority: priority,
		RespChan: respChan,
	}
	
	// Push request into priority queue (thread-safe)
	d.mu.Lock()
	heap.Push(d.queue,req)
	d.mu.Unlock()
	
	// Wait for allocation result from worker
	return <-respChan
}

// Run is a long-running worker that continuously processes queued requests.
//
// Behavior:
// - Always picks the highest priority request
// - Delegates allocation to Allocator
// - Sends result back via response channel
//
// WARNING:
// - This loop is intentionally infinite
// - Avoid busy-waiting in production (can be improved using condition variables or channels)
func (d *Dispatcher) Run() {
	for {
		d.mu.Lock()

		// If no requests are available, release lock and retry
		if d.queue.Len() == 0 {
			d.mu.Unlock()
			continue
		}

		// Pop highest-priority request
		req := heap.Pop(d.queue).(Request)
		d.mu.Unlock()
		
		v := req.Vehicle

		// Allocate parking slot
		level,slot,err := d.allocator.Allocate(v)
		
		// Send result back to requester
		req.RespChan <- Result{
			Level: level,
			Slot: slot,
			Err: err,
		}
	}
}

// getPriority maps customer type to priority value.
//
// Priority Rules:
// - Emergency → Highest priority (1)
// - VIP       → Medium priority (2)
// - Regular   → Lowest priority (3)
//
// Lower number = higher priority
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