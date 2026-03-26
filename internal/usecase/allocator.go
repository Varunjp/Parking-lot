package usecase

import (
	"math"
	"parking-lot/internal/domain"
	"time"
)

type Allocator struct {
	parkingRepo domain.ParkingRepository
	vehicleRepo domain.VehicleRepository
}

func NewAllocator(p domain.ParkingRepository,v domain.VehicleRepository) *Allocator {
	return &Allocator{p,v}
}

func (a *Allocator) Allocate(vehicle domain.Vehicle)(int,int,error) {

	// Check re-entry within 1 hour
	if rec,ok := a.vehicleRepo.GetLastEntry(vehicle.ID); ok {
		if time.Now().Unix() - rec.LastEntryTime < 3600 {
			return 0,0,domain.ErrReEntry
		}
	}

	var bestLevel *domain.Level
	min := math.MaxInt

	for _,lvl := range a.parkingRepo.GetLevels() {
		pool := getPool(lvl,vehicle.Type)
		if pool == nil {
			continue
		}

		if len(pool.FreeSlots) > 0 {
			if len(pool.Occupied) < min {
				min = len(pool.Occupied)
				bestLevel = lvl 
			}
		}
	}

	if bestLevel == nil {
		return 0,0,domain.ErrParkingFull
	}

	pool := getPool(bestLevel,vehicle.Type)
	slot := pool.FreeSlots[len(pool.FreeSlots)-1]
	pool.FreeSlots = pool.FreeSlots[:len(pool.FreeSlots)-1]
	pool.Occupied[slot] = true

	a.vehicleRepo.SaveEntry(vehicle.ID,domain.VehicleRecord{
		LastEntryTime: time.Now().Unix(),
		LastLevel: bestLevel.ID,
	})

	return bestLevel.ID,slot,nil 
}

func getPool(l *domain.Level,t domain.VehicleType) *domain.SlotPool {
	switch t {
	case domain.Small:
		return l.SmallSlots
	case domain.Medium:
		return l.MediumSlots
	case domain.Large:
		return l.LargeSlots
	}
	return nil 
}