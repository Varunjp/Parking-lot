package memory

import "parking-lot/internal/domain"

type ParkingRepo struct {
	ParkingLot  *domain.ParkingLot
}

func (r *ParkingRepo) GetLevels()[]*domain.Level {
	return r.ParkingLot.Levels
}