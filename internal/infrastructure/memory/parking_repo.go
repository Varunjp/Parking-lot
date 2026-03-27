package memory

import "parking-lot/internal/domain"

// ParkingRepo is an in-memory implementation of the ParkingRepository interface.
//
// It holds a reference to the ParkingLot aggregate and provides read access
// to parking-related data such as levels and slots.
//
// NOTE:
// - This implementation is intended for testing or lightweight use cases.
// - It is NOT safe for concurrent access unless explicitly synchronized.
type ParkingRepo struct {
	// ParkingLot represents the root aggregate containing all parking levels.
	// It should be initialized before using the repository.
	ParkingLot  *domain.ParkingLot
}

// GetLevels returns all parking levels available in the parking lot.
//
// Returns:
// - []*domain.Level: A slice of pointers to Level objects.
//
// Behavior:
// - If ParkingLot is nil, it returns nil (caller should handle this case).
// - The returned slice is not a deep copy; modifying it may affect internal state.
//
// Recommendation:
// - Treat the returned data as read-only unless mutation is intended.
func (r *ParkingRepo) GetLevels()[]*domain.Level {
	return r.ParkingLot.Levels
}