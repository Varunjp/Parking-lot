package domain

// ParkingRepository defines the contract for accessing and managing
// parking lot state (levels, slots, availability).
//
// This abstraction allows the service layer to remain independent of
// how parking data is stored (in-memory, DB, cache, etc.).
type ParkingRepository interface {
	// GetLevels returns all parking levels in the system.
	//
	// Each Level contains its slot pools (e.g., small, medium, large)
	// along with current availability and occupancy details.
	//
	// Expected behavior:
	// - Always returns the latest state of parking levels
	// - Should not return nil; return an empty slice if no levels exist
	GetLevels() []*Level
}

// VehicleRepository defines the contract for tracking vehicle
// entry history and enforcing re-entry constraints.
//
// This is typically used for:
// - Preventing immediate re-entry to the same level
// - Maintaining parking history for auditing or analytics
type VehicleRepository interface {
	
	// GetLastEntry retrieves the last known parking record(s)
	// for a given vehicle.
	//
	// Parameters:
	// - vehicleID: unique identifier of the vehicle (e.g., plate number)
	//
	// Returns:
	// - map[levelID]timestamp:
	//     Key   -> Level ID where the vehicle was previously parked
	//     Value -> Unix timestamp of the last entry/exit
	// - bool:
	//     true  -> history exists for this vehicle
	//     false -> no prior record found
	//
	// Notes:
	// - Multiple entries may exist if the vehicle has used different levels
	// - Used for enforcing rules like "no re-entry to same level within X time"
	GetLastEntry(vehicleID string) (map[int]int64, bool)

	// SaveEntry records a new parking entry for a vehicle.
	//
	// Parameters:
	// - vehicleID: unique identifier of the vehicle
	// - levelID:   level where the vehicle is parked
	// - ts:        Unix timestamp representing the entry time
	//
	// Expected behavior:
	// - Should update or insert the latest entry for the vehicle
	// - Must ensure consistency if multiple entries are tracked
	// - Can be extended to store exit time or duration if needed
	SaveEntry(vehicleID string, levelID int, ts int64)

	GetActive(vehicleID string)(ActiveParking,bool)
	SaveActive(vehicleID string,levelID int,slotID int,slotType *SlotPool,ts int64)
	RemoveActive(vechileID string)
}
