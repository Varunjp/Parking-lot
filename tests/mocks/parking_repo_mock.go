package mock

import "parking-lot/internal/domain"

type MockParkingRepo struct {
	Levels []*domain.Level
}

func (m *MockParkingRepo) GetLevels() []*domain.Level {
	return m.Levels
}