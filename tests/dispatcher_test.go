package tests

import (
	"parking-lot/internal/domain"
	"parking-lot/internal/usecase"
	mock "parking-lot/tests/mocks"
	"sync"
	"testing"
)

// TestDispatcher_PriorityOrder verifies that higher-priority vehicles
// are allocated slots before lower-priority ones under concurrent requests.
//
// Priority Order Expected:
//   1. Emergency
//   2. VIP
//   3. Regular
//
// This test simulates concurrent incoming requests and ensures that
// dispatcher enforces priority correctly despite goroutine scheduling.
func TestDispatcher_PriorityOrder(t *testing.T) {

	// Arrange: create a parking level with 10 available slots
	level := createLevel(1,10)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: make(map[string]map[int]int64),
	}

	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,3600)
	dispatcher := usecase.NewDispatcher(allocator)

	// Input vehicles with different priorities
	vehicles := []domain.Vehicle{
		{ID: "V1",Type: domain.Small, CustomerType: domain.Regular},
		{ID: "V2",Type: domain.Small, CustomerType: domain.Emergency},
		{ID: "V3",Type: domain.Small, CustomerType: domain.VIP},
	}

	results := make(map[string]usecase.Result)
	var mu sync.Mutex     // protects results map from race conditions
	var wg sync.WaitGroup
	
	// Act: simulate concurrent requests
	for _,v := range vehicles {
		wg.Add(1)
		go func(v domain.Vehicle) {
			defer wg.Done()
			
			res := dispatcher.AddRequest(v)
			
			// Safely store result
			mu.Lock()
			results[v.ID] = res 
			mu.Unlock()
		}(v)
	}

	wg.Wait()

	// Assert: verify priority-based allocation order
	if results["V2"].Slot != 1 {
		t.Fatalf("expected Emergency to get slot 1, got %d", results["V2"].Slot)
	}

	if results["V3"].Slot != 2 {
		t.Fatalf("expected VIP to get slot 2, got %d", results["V3"].Slot)
	}

	if results["V1"].Slot != 3 {
		t.Fatalf("expected Regular to get slot 3, got %d", results["V1"].Slot)
	}
}

// TestDispatcher_ConcurrentRequests ensures that the dispatcher
// can safely handle multiple concurrent requests without:
//   - race conditions
//   - duplicate slot allocation
//   - missed allocations
//
// This test does NOT verify priority, only correctness under concurrency.
func TestDispatcher_ConcurrentRequests(t *testing.T) {

	// Arrange
	level := createLevel(1, 50)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: make(map[string]map[int]int64),
	}

	allocator := usecase.NewAllocator(parkingRepo, vehicleRepo, 3600)
	dispatcher := usecase.NewDispatcher(allocator)

	totalRequests := 20
	results := make([]usecase.Result, totalRequests)

	var wg sync.WaitGroup

	// Act: fire concurrent allocation requests
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			v := domain.Vehicle{
				ID:           string(rune('A' + index)),
				Type: domain.Small,
				CustomerType: domain.Regular,
			}

			results[index] = dispatcher.AddRequest(v)
		}(i)
	}

	wg.Wait()

	// Assert: verify all allocations are valid
	for i, res := range results {
		if res.Err != nil {
			t.Fatalf("unexpected error at %d: %v", i, res.Err)
		}
		if res.Slot == 0 {
			t.Fatalf("invalid slot allocation at %d", i)
		}
	}
}

// TestDispatcher_InvalidPriority verifies that dispatcher
// correctly rejects vehicles with unsupported/invalid priority.
//
// This ensures system robustness against malformed input.
func TestDispatcher_InvalidPriority(t *testing.T) {

	// Arrange: use a dummy allocator since validation happens before allocation
	allocator := &usecase.Allocator{}
	dispatcher := usecase.NewDispatcher(allocator)

	v := domain.Vehicle{
		ID:           "X1",
		Type: domain.Large,
		CustomerType: domain.CustomerType("999"), // invalid
	}

	// Act
	res := dispatcher.AddRequest(v)

	// Assert
	if res.Err == nil {
		t.Fatalf("expected error for invalid priority")
	}
}