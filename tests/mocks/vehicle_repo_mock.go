package mock

type MockVehicleRepo struct {
	Records map[string]map[int]int64
}

func (m *MockVehicleRepo) GetLastEntry(vehicleID string) (map[int]int64, bool) {
	rec, ok := m.Records[vehicleID]
	return rec, ok
}

func (m *MockVehicleRepo) SaveEntry(vehicleID string, levelID int, ts int64) {
	if m.Records == nil {
		m.Records = make(map[string]map[int]int64)
	}

	if m.Records[vehicleID] == nil {
		m.Records[vehicleID] = make(map[int]int64)
	}

	m.Records[vehicleID][levelID] = ts
}