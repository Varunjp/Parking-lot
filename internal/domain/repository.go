package domain

type ParkingRepository interface {
	GetLevels() []*Level
}

type VehicleRepository interface {
	GetLastEntry(vehicleID string) (VehicleRecord, bool)
	SaveEntry(vehicleID string, record VehicleRecord)
}

type VehicleRecord struct {
	LastEntryTime int64
	LastLevel     int
}