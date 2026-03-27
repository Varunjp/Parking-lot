package main

import (
	"log"
	"net/http"
	"parking-lot/internal/config"
	httpClient "parking-lot/internal/delivery/http"
	"parking-lot/internal/domain"
	"parking-lot/internal/infrastructure/memory"
	"parking-lot/internal/usecase"
)

// main is the application entry point.
// It wires together configuration, repositories, use cases, and HTTP handlers,
// then starts the HTTP server.
func main() {

	// Load application configuration (env / defaults)
	cfg := config.Load()

	// Initialize in-memory repositories
	vehicleRepo := memory.NewVehicleRepo()

	// Initialize parking lot structure (levels + slot pools)
	parkingLot, err := initLevels(cfg)
	if err != nil {
		log.Fatal("failed to initialize parking levels: ", err)
	}

	parkingRepo := &memory.ParkingRepo{
		ParkingLot: parkingLot,
	}

	// Initialize core business logic (use cases)
	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,cfg.ReEntrySeconds)
	dispatcher := usecase.NewDispatcher(allocator)

	// Setup HTTP handler layer
	handler := httpClient.NewHandler(dispatcher)

	// Register routes
	http.HandleFunc("/park",handler.Park)

	// Start HTTP server
	log.Println("Server running on :"+cfg.HttpPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HttpPort,nil))
}

// initLevels constructs the parking lot with multiple levels.
// Each level contains separate slot pools for small, medium, and large vehicles.
//
// WHY this exists:
// - Keeps initialization logic separate from main()
// - Makes testing and future extension easier (e.g., DB-backed config)
// - Central place to enforce parking constraints
func initLevels(cfg *config.Config) (*domain.ParkingLot,error) {

	var levels []*domain.Level

	for i:=1; i <= cfg.ParkingLevels; i++ {

		// Initialize slot pools per vehicle size
		// FreeSlots acts as a stack/queue for O(1) allocation
		smallPool := &domain.SlotPool{
			FreeSlots: generateSlot(cfg.SmallSlotsPerLevel),
			Occupied: make(map[int]bool),
		}

		mediumPool := &domain.SlotPool{
			FreeSlots: generateSlot(cfg.MediumSlotsPerLevel),
			Occupied: make(map[int]bool),
		}

		largePool := &domain.SlotPool{
			FreeSlots: generateSlot(cfg.LargeSlotsPerLevel),
			Occupied: make(map[int]bool),
		}

		// Create level with validation handled inside domain
		level,err := domain.NewLevel(i,smallPool,mediumPool,largePool)
		if err != nil {
			return nil,err 
		}

		levels = append(levels, level)
	}

	return &domain.ParkingLot{Levels: levels},nil 
}

// generateSlots creates a sequential list of slot IDs starting from 1.
//
// WHY:
// - Ensures deterministic slot numbering
// - Keeps allocation logic simple (pop from FreeSlots)
// - Avoids runtime generation overhead during allocation
//
// Example:
// n = 3 → [1, 2, 3]
func generateSlot(n int) []int {
	slots := make([]int,n)
	for i := 0; i < n; i++ {
		slots[i] = i + 1
	}
	return slots
}