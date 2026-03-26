package domain

type ParkingRepository interface {
	GetLevels() []*Level
}

type VehicleRepository interface {
	GetLastEntry(vehicleID string) (map[int]int64, bool)
	SaveEntry(vehicleID string, levelID int, ts int64)
}

type VehicleRecord struct {
	LastEntryTime int64
	LastLevel     int
}