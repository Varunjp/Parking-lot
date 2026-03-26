package memory

import (
	"sync"
)

type VehicleRepo struct {
	mu sync.RWMutex
	data map[string]map[int]int64
}

func NewVehicleRepo() *VehicleRepo {
	return &VehicleRepo{
		data: make(map[string]map[int]int64),
	}
}

func (r *VehicleRepo) GetLastEntry(id string)(map[int]int64,bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res,ok := r.data[id]
	return res,ok 
}

func (r *VehicleRepo) SaveEntry(id string, levelID int,ts int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _,ok := r.data[id]; !ok {
		r.data[id] = make(map[int]int64)
	}
	r.data[id][levelID] = ts
}