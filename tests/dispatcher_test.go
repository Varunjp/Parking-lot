package tests

import (
	"parking-lot/internal/domain"
	"parking-lot/internal/usecase"
	mock "parking-lot/tests/mocks"
	"sync"
	"testing"
	"time"
)

func TestDispatcher_PriorityOrder(t *testing.T) {

	level := createLevel(1,10)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: make(map[string]map[int]int64),
	}

	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,3600)
	dispatcher := usecase.NewDispatcher(allocator)

	vehicles := []domain.Vehicle{
		{ID: "V1",Type: domain.Small, CustomerType: domain.Regular},
		{ID: "V2",Type: domain.Small, CustomerType: domain.Emergency},
		{ID: "V3",Type: domain.Small, CustomerType: domain.VIP},
	}

	results := make(map[string]usecase.Result)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for _,v := range vehicles {
		wg.Add(1)
		go func(v domain.Vehicle) {
			defer wg.Done()
			res := dispatcher.AddRequest(v)
			
			mu.Lock()
			results[v.ID] = res 
			mu.Unlock()
		}(v)
	}

	wg.Wait()

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

func TestDispatcher_ConcurrentRequests(t *testing.T) {

	level := createLevel(1, 50)

	parkingRepo := &mock.MockParkingRepo{
		Levels: []*domain.Level{level},
	}

	vehicleRepo := &mock.MockVehicleRepo{
		Records: make(map[string]map[int]int64),
	}

	allocator := usecase.NewAllocator(parkingRepo, vehicleRepo, 3600)
	dispatcher := usecase.NewDispatcher(allocator)

	total := 20
	results := make([]usecase.Result, total)

	for i := 0; i < total; i++ {
		i := i
		go func() {
			v := domain.Vehicle{
				ID:           string(rune('A' + i)),
				Type: domain.Small,
				CustomerType: domain.Regular,
			}
			results[i] = dispatcher.AddRequest(v)
		}()
	}

	time.Sleep(100 * time.Millisecond)

	for i, res := range results {
		if res.Err != nil {
			t.Fatalf("unexpected error at %d: %v", i, res.Err)
		}
		if res.Slot == 0 {
			t.Fatalf("invalid slot allocation at %d", i)
		}
	}
}

func TestDispatcher_InvalidPriority(t *testing.T) {

	allocator := &usecase.Allocator{} // dummy
	dispatcher := usecase.NewDispatcher(allocator)

	v := domain.Vehicle{
		ID:           "X1",
		Type: domain.Large,
		CustomerType: domain.CustomerType("999"), // invalid
	}

	res := dispatcher.AddRequest(v)

	if res.Err == nil {
		t.Fatalf("expected error for invalid priority")
	}
}