package usecase

import (
	"math"
	"parking-lot/internal/domain"
	"time"
)

type Allocator struct {
	parkingRepo domain.ParkingRepository
	vehicleRepo domain.VehicleRepository
	reentrySec  int64
}

func NewAllocator(p domain.ParkingRepository,v domain.VehicleRepository,reentry int64) *Allocator {
	return &Allocator{p,v,reentry}
}

func (a *Allocator) Allocate(vehicle domain.Vehicle)(int,int,error) {

	rec,hasHistory := a.vehicleRepo.GetLastEntry(vehicle.ID)
	now := time.Now().Unix()

	var bestLevel *domain.Level
	min := math.MaxInt

	for _,lvl := range a.parkingRepo.GetLevels() {
		pool,err := getPool(lvl,vehicle.Type)
		if err != nil {
			return 0,0,err 
		}
		if pool == nil || len(pool.FreeSlots) == 0{
			continue
		}

		if hasHistory {
			if lastTime,ok := rec[lvl.ID]; ok {
				if now-lastTime < a.reentrySec {
					continue
				}
			}
		}

		if len(pool.FreeSlots) > 0 {
			if len(pool.Occupied) < min {
				min = len(pool.Occupied)
				bestLevel = lvl 
			}
		}
	}

	if bestLevel == nil {
		if hasHistory {
			return 0,0,domain.ErrReEntry
		}
		return 0,0,domain.ErrParkingFull
	}

	pool,err := getPool(bestLevel,vehicle.Type)
	if err != nil {
		return 0,0,err 
	}
	slot := pool.FreeSlots[len(pool.FreeSlots)-1]
	pool.FreeSlots = pool.FreeSlots[:len(pool.FreeSlots)-1]
	pool.Occupied[slot] = true

	a.vehicleRepo.SaveEntry(vehicle.ID,bestLevel.ID,now)

	return bestLevel.ID,slot,nil 
}

func getPool(l *domain.Level,t domain.VehicleType) (*domain.SlotPool,error) {
	switch t {
	case domain.Small:
		return l.SmallSlots,nil
	case domain.Medium:
		return l.MediumSlots,nil
	case domain.Large:
		return l.LargeSlots,nil
	}
	return nil,domain.ErrInvalidType
}