package mock

import "parking-lot/internal/domain"

// MockParkingRepo is an in-memory implementation of the ParkingRepository interface.
//
// It is primarily used for unit testing to simulate parking lot data
// without relying on a real database or external service.
//
// This mock allows controlled setup of parking levels and predictable behavior
// during tests.
type MockParkingRepo struct {
	// Levels represents the list of parking levels available in the system.
	// This can be pre-populated in tests to simulate different parking scenarios.
	Levels []*domain.Level
}

// GetLevels returns all available parking levels.
//
// This method satisfies the ParkingRepository interface and provides
// a deterministic response for testing purposes.
//
// Returns:
//   - []*domain.Level: A slice of parking levels configured in the mock.
func (m *MockParkingRepo) GetLevels() []*domain.Level {
	return m.Levels
}