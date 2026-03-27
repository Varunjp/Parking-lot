package usecase

import (
	"container/heap"
	"parking-lot/internal/domain"
)

// Request represents a parking allocation request.
//
// It encapsulates:
// - Vehicle details
// - Priority used for scheduling (lower value = higher priority)
// - A response channel to asynchronously return the allocation result
type Request struct {
	Vehicle  domain.Vehicle // Incoming vehicle requesting a slot
	Priority int			// Scheduling priority (min-heap: lower = higher priority)
	RespChan chan Result 	// Channel to send allocation result back to caller
}

// PQ (Priority Queue) implements a min-heap of Requests based on Priority.
//
// This is used by the dispatcher to ensure that higher-priority vehicles
// (e.g., emergency) are processed before others.
//
// NOTE:
// The heap is ordered such that the smallest Priority value is popped first.
type PQ []Request

// Len returns the number of elements in the priority queue.
//
// Required by heap.Interface.
func (pq PQ) Len() int { 
	return len(pq) 
}

// Less determines the ordering of elements in the heap.
// Returns true if element at index i should come before element at index j.
// We use a min-heap based on Priority:
// - Lower Priority value = higher importance
// - Example: Priority 1 (Emergency) > Priority 2 (Normal)
func (pq PQ)Less(i, j int) bool { 
	return pq[i].Priority < pq[j].Priority 
}

// Swap exchanges elements at indices i and j.
// Required by heap.Interface.
func (pq PQ) Swap(i, j int){ 
	pq[i], pq[j] = pq[j], pq[i] 
}

// Push adds a new element to the priority queue.
//
// This method is called internally by heap.Push().
// It appends the new Request to the underlying slice.
//
// IMPORTANT:
// The input type is interface{} due to heap.Interface contract,
// so we must type assert it to Request.
func (pq *PQ) Push(x interface{}) {
	*pq = append(*pq, x.(Request))
}

// Pop removes and returns the highest-priority element from the queue.
//
// This method is called internally by heap.Pop().
// It removes the last element (after heap adjustment).
//
// IMPORTANT:
// Always ensure the queue is not empty before calling Pop externally.
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)

	if n == 0 {
		return nil // Defensive safety (should not happen if used correctly)
	}
	// Extract last element
	req := old[n-1]

	// Shrink slice
	*pq = old[:n-1]

	return req
}

// Init initializes the priority queue heap structure.
//
// Must be called before using heap operations if the PQ
// is not already heapified.
//
// Example usage:
//   pq := &PQ{}
//   pq.Init()
func (pq *PQ) Init() {
	heap.Init(pq)
}