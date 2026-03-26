package memory

import "parking-lot/internal/domain"

type ParkingRepo struct {
	Levels []*domain.Level
}

func (r *ParkingRepo) GetLevels()[]*domain.Level {
	return r.Levels
}