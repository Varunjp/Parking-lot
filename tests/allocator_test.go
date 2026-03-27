package tests

import (
	"parking-lot/internal/domain"
	"parking-lot/internal/usecase"
	mock "parking-lot/tests/mocks"
	"testing"
	"time"
)

/*
TestAllocate_Success verifies that a vehicle is successfully allocated
to an available slot when:
- Slots are available across levels
- No re-entry restriction applies

Expected:
- No error
- Valid levelID and slotID (> 0)
*/
func TestAllocate_Success(t *testing.T) {
	// Arrange: Create parking levels with available slots
	level1 := createLevel(1,2)
	level2 := createLevel(2,2)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level1,level2},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		// No previous entry → no re-entry restriction
		Records: make(map[string]map[int]int64),
	}

	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,3600)

	vehicle := domain.Vehicle{
		ID: "KL01",
		Type: domain.Small,
	}

	// Act
	levelID,slotID,err := allocator.Allocate(vehicle)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v",err)
	}

	if levelID == 0 || slotID == 0 {
		t.Fatalf("expected valid allocation, got level=%d slot=%d",levelID,slotID)
	}
}

/*
TestAllocate_ReEntryBlocked verifies that a vehicle is blocked
from re-entering the same level within the restricted time window.

Scenario:
- Vehicle previously parked in level 1
- Re-entry attempted within restriction time (3600 seconds)

Expected:
- Allocation should fail with an error
*/
func TestAllocate_ReEntryBlocked(t *testing.T) {
	// Arrange: Create a level with available slots
	level := createLevel(1,2)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: map[string]map[int]int64{
			"KL01":{
				// Simulate recent exit (within restriction window)
				1: time.Now().Unix(),
			},
		},
	}

	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,3600)

	vehicle := domain.Vehicle {
		ID: "KL01",
		Type: domain.Small,
	}

	// Act
	_,_,err := allocator.Allocate(vehicle)

	// Assert
	if err == nil {
		t.Fatalf("expected re-entry error, got nil")
	}
}

/*
TestAllocate_NoSlots verifies behavior when no slots are available.

Scenario:
- Level exists but has zero free slots

Expected:
- Allocation should fail with an error
*/
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

	// Act
	_,_,err := allocator.Allocate(vehicle)

	// Assert
	if err == nil {
		t.Fatalf("expected error when no slots available")
	}
}

/*
TestAllocate_SelectBestLevel verifies that the allocator selects
the most optimal level based on slot availability.

Scenario:
- Level 1 has no free slots
- Level 2 has available slots

Expected:
- Allocation should happen in level 2
*/
func TestAllocate_SelectBestLevel(t *testing.T) {
	// Arrange
	level1 := createLevel(1, 1)
	level2 := createLevel(2, 3)

	// Simulate level1 being fully occupied
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

	// Act
	levelID, _, _ := allocator.Allocate(vehicle)

	// Assert
	if levelID != 2 {
		t.Fatalf("expected allocation in level 2, got %d", levelID)
	}
}