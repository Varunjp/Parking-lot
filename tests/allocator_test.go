package tests

import (
	"parking-lot/internal/domain"
	"parking-lot/internal/usecase"
	mock "parking-lot/tests/mocks"
	"testing"
	"time"
)

func TestAllocate_Success(t *testing.T) {

	level1 := createLevel(1,2)
	level2 := createLevel(2,2)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level1,level2},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: make(map[string]map[int]int64),
	}

	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,3600)

	vehicle := domain.Vehicle{
		ID: "KL01",
		Type: domain.Small,
	}

	levelID,slotID,err := allocator.Allocate(vehicle)

	if err != nil {
		t.Fatalf("expected no error, got %v",err)
	}

	if levelID == 0 || slotID == 0 {
		t.Fatalf("expected valid allocation, got level=%d slot=%d",levelID,slotID)
	}
}

func TestAllocate_ReEntryBlocked(t *testing.T) {

	level := createLevel(1,2)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: map[string]map[int]int64{
			"KL01":{
				1: time.Now().Unix(),
			},
		},
	}

	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,3600)

	vehicle := domain.Vehicle {
		ID: "KL01",
		Type: domain.Small,
	}

	_,_,err := allocator.Allocate(vehicle)

	if err == nil {
		t.Fatalf("expected re-entry error, got nil")
	}
}

func TestAllocate_NoSlots(t *testing.T) {

	level := createLevel(1,0)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: make(map[string]map[int]int64),
	}

	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,3600)

	vehicle := domain.Vehicle {
		ID: "KL01",
		Type: domain.Small,
	}

	_,_,err := allocator.Allocate(vehicle)

	if err == nil {
		t.Fatalf("expected error when no slots available")
	}
}

func TestAllocate_SelectBestLevel(t *testing.T) {

	level1 := createLevel(1, 1)
	level2 := createLevel(2, 3)

	level1.SmallSlots.FreeSlots = []int{}

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level1, level2},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: make(map[string]map[int]int64),
	}

	allocator := usecase.NewAllocator(parkingRepo, vehicleRepo, 3600)

	vehicle := domain.Vehicle{
		ID:   "KL01",
		Type: domain.Small,
	}

	levelID, _, _ := allocator.Allocate(vehicle)

	if levelID != 2 {
		t.Fatalf("expected allocation in level 2, got %d", levelID)
	}
}