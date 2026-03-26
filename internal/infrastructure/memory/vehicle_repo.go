package memory

import (
	"parking-lot/internal/domain"
	"sync"
)

type VehicleRepo struct {
	mu sync.RWMutex
	data map[string]domain.VehicleRecord
}

func NewVehicleRepo() *VehicleRepo {
	return &VehicleRepo{
		data: make(map[string]domain.VehicleRecord),
	}
}

func (r *VehicleRepo) GetLastEntry(id string)(domain.VehicleRecord,bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res,ok := r.data[id]
	return res,ok 
}

func (r *VehicleRepo) SaveEntry(id string,rec domain.VehicleRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[id] = rec
}