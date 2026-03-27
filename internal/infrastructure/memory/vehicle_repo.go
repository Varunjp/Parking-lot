package memory

import (
	"sync"
)

// VehicleRepo provides a thread-safe in-memory storage
// for tracking vehicle entry timestamps per parking level.
//
// Structure:
//   vehicleID -> (levelID -> lastEntryTimestamp)
//
// Example:
//   "KL-01-1234" -> {
//       1: 1711540000,
//       2: 1711543600,
//   }
//
// This is primarily used for:
//   - Enforcing re-entry restrictions
//   - Tracking last known parking activity
type VehicleRepo struct {
	mu sync.RWMutex
	data map[string]map[int]int64
}

// NewVehicleRepo initializes and returns a new VehicleRepo instance.
func NewVehicleRepo() *VehicleRepo {
	return &VehicleRepo{
		data: make(map[string]map[int]int64),
	}
}

// GetLastEntry returns a COPY of the last entry timestamps for a given vehicle.
//
// Returns:
//   - map[levelID]timestamp → last entry time per level
//   - bool → false if vehicle has no history
//
// NOTE:
// A copy is returned to prevent accidental mutation of internal state.
func (r *VehicleRepo) GetLastEntry(id string)(map[int]int64,bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	levelData,exists := r.data[id]
	if !exists {
		return nil,false
	}

	// Create a copy to avoid exposing internal map
	result := make(map[int]int64, len(levelData))
	for levelID, ts := range levelData {
		result[levelID] = ts
	}

	return result,true 
}

// SaveEntry records the latest entry timestamp for a vehicle at a specific level.
//
// If the vehicle does not exist, it initializes storage for it.
//
// Parameters:
//   - vehicleID: unique identifier of the vehicle
//   - levelID: parking level where vehicle is entering
//   - ts: timestamp of entry (Unix time)
func (r *VehicleRepo) SaveEntry(vehicleid string, levelID int,ts int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Initialize vehicle record if not present
	if _,ok := r.data[vehicleid]; !ok {
		r.data[vehicleid] = make(map[int]int64)
	}

	r.data[vehicleid][levelID] = ts
}