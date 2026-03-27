package usecase

import (
	"math"
	"parking-lot/internal/domain"
	"time"
)

// Allocator is responsible for assigning parking slots to vehicles
// based on availability, load balancing, and re-entry constraints.

type Allocator struct {
	parkingRepo domain.ParkingRepository
	vehicleRepo domain.VehicleRepository
	reentrySec  int64 // Restriction window for re-entering same level
}

// NewAllocator initializes a new Allocator instance with required dependencies.
func NewAllocator(p domain.ParkingRepository,v domain.VehicleRepository,reentry int64) *Allocator {
	return &Allocator{p,v,reentry}
}

// Allocate assigns a parking slot to the given vehicle.
//
// Allocation strategy:
// 1. Select levels that support the vehicle type and have free slots.
// 2. Apply re-entry restriction (prevent same level within time window).
// 3. Choose the level with the least occupied slots (load balancing).
// 4. Allocate slot and persist entry.
//
// Returns:
// - levelID: ID of allocated level
// - slotID: allocated slot number
// - error: if allocation fails (invalid type, full, or restricted)

func (a *Allocator) Allocate(vehicle domain.Vehicle)(int,int,error) {

	now := time.Now().Unix()

	// Fetch vehicle history for re-entry validation
	lastEntryByLevel,hasHistory := a.vehicleRepo.GetLastEntry(vehicle.ID)
	
	var (
		selectedLevel *domain.Level
		selectedPool  *domain.SlotPool
		minOccupied   = math.MaxInt
		blockedLevels = 0
		totalLevels   = 0
	)

	// Iterate through all levels to find the best candidate
	for _,level := range a.parkingRepo.GetLevels() {
		totalLevels++
		
		pool,ok := level.GetPool(vehicle.Type)
		if !ok {
			return 0,0,domain.ErrInvalidType
		}

		// Skip levels with no available slots
		if pool.AvailableCount() == 0{
			continue
		}

		// Enforce re-entry restriction per level
		if a.isReentryBlocked(hasHistory, lastEntryByLevel, level.ID, now) {
			blockedLevels++
			continue
		}

		// Select level with minimum occupancy (load balancing)
		if len(pool.Occupied) < minOccupied {
			minOccupied = len(pool.Occupied)
			selectedLevel = level 
			selectedPool = pool
		}
	}

	// Handle case when no suitable level found
	if selectedLevel == nil {
		return a.handleAllocationFailure(hasHistory,blockedLevels,totalLevels)
	}

	// Allocate slot from selected pool
	slot,err := selectedPool.Allocate()
	if err != nil {
		return 0,0,err 
	}

	// Persist vehicle entry for future re-entry checks
	a.vehicleRepo.SaveEntry(vehicle.ID,selectedLevel.ID,now)

	return selectedLevel.ID,slot,nil 
}

//
// ------------------- Helper Methods -------------------
//

// isReentryBlocked checks whether the vehicle is restricted from entering
// the given level due to recent exit within the configured time window.
func (a *Allocator) isReentryBlocked(
	hasHistory bool,
	lastEntry map[int]int64,
	levelID int,
	currentTime int64,
) bool {

	if !hasHistory {
		return false
	}

	lastTime, exists := lastEntry[levelID]
	if !exists {
		return false
	}

	// Block if within re-entry restriction window
	return currentTime-lastTime < a.reentrySec
}

// handleAllocationFailure determines the appropriate error when allocation fails.
func (a *Allocator) handleAllocationFailure(
	hasHistory bool,
	blockedLevels int,
	totalLevels int,
) (int, int, error) {

	// All levels blocked due to re-entry restriction
	if hasHistory && blockedLevels == totalLevels {
		return 0, 0, domain.ErrReEntry
	}

	// No slots available across all levels
	return 0, 0, domain.ErrParkingFull
}