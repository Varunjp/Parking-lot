package usecase

import "container/heap"

type Request struct {
	Vehicle  interface{}
	Priority int
}

type PQ []Request

func (pq PQ) Len() int           { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }
func (pq PQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PQ) Push(x interface{}) {
	*pq = append(*pq, x.(Request))
}

func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func (pq *PQ) Init() {
	heap.Init(pq)
}