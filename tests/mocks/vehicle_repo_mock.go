package mock

// MockVehicleRepo is an in-memory implementation of the VehicleRepository.
// It is primarily used for unit testing and simulates persistence behavior
// without requiring a real database.
//
// Data Structure:
// - Outer map key   -> vehicleID (string)
// - Inner map key   -> levelID (int)
// - Inner map value -> timestamp of entry (int64)
//
// Example:
// Records = {
//   "KL-01-1234": {
//       1: 1710000000, // level 1 entry timestamp
//       2: 1710003600, // level 2 entry timestamp
//   },
// }
type MockVehicleRepo struct {
	Records map[string]map[int]int64
}

// GetLastEntry retrieves the last recorded entry timestamps for a given vehicle.
//
// Parameters:
// - vehicleID: unique identifier of the vehicle
//
// Returns:
// - map[int]int64: map of levelID to entry timestamp
// - bool: true if records exist for the vehicle, false otherwise
//
// Notes:
// - This does NOT return a single "last" entry, but all recorded level entries.
// - Caller is responsible for determining the most recent entry if needed.
func (m *MockVehicleRepo) GetLastEntry(vehicleID string) (map[int]int64, bool) {
	rec, ok := m.Records[vehicleID]
	return rec, ok
}

// SaveEntry stores or updates the entry timestamp for a vehicle at a specific level.
//
// Parameters:
// - vehicleID: unique identifier of the vehicle
// - levelID: parking level where the vehicle is entering
// - ts: Unix timestamp representing entry time
//
// Behavior:
// - Initializes internal maps if they are nil
// - Overwrites existing timestamp if the vehicle already has an entry for the same level
//
// Example:
// SaveEntry("KL-01-1234", 1, 1710000000)
func (m *MockVehicleRepo) SaveEntry(vehicleID string, levelID int, ts int64) {
	if m.Records == nil {
		m.Records = make(map[string]map[int]int64)
	}

	if m.Records[vehicleID] == nil {
		m.Records[vehicleID] = make(map[int]int64)
	}

	m.Records[vehicleID][levelID] = ts
}