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
	var bestPool *domain.SlotPool
	min := math.MaxInt
	blockedLevels := 0
	totalLevels := 0

	for _,lvl := range a.parkingRepo.GetLevels() {
		totalLevels++
		
		pool,ok := lvl.GetPool(vehicle.Type)
		if !ok {
			return 0,0,domain.ErrInvalidType
		}
		if pool == nil || pool.AvailableCount() == 0{
			continue
		}
		if hasHistory {
			if lastTime,ok := rec[lvl.ID]; ok {
				if now-lastTime < a.reentrySec {
					blockedLevels++
					continue
				}
			}
		}
		if len(pool.Occupied) < min {
			min = len(pool.Occupied)
			bestLevel = lvl 
			bestPool = pool
		}
	}

	if bestLevel == nil {
		if hasHistory && blockedLevels == totalLevels {
			return 0,0,domain.ErrReEntry
		}
		return 0,0,domain.ErrParkingFull
	}

	slot,err := bestPool.Allocate()
	if err != nil {
		return 0,0,err 
	}

	a.vehicleRepo.SaveEntry(vehicle.ID,bestLevel.ID,now)

	return bestLevel.ID,slot,nil 
}
