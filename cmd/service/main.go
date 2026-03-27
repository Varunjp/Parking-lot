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

func main() {

	cfg := config.Load()

	vehicleRepo := memory.NewVehicleRepo()
	levels,err := initLevels(cfg)

	if err != nil {
		log.Fatal(err)
	}

	parkingRepo := &memory.ParkingRepo{
		ParkingLot: levels,
	}

	allocator := usecase.NewAllocator(parkingRepo,vehicleRepo,cfg.ReEntrySeconds)
	dispatcher := usecase.NewDispatcher(allocator)

	go func() {
		for {
			dispatcher.Process()
		}
	}()

	handler := httpClient.NewHandler(dispatcher)
	http.HandleFunc("/park",handler.Park)

	log.Println("Server running on :"+cfg.HttpPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HttpPort,nil))
}

func initLevels(cfg *config.Config) (*domain.ParkingLot,error) {

	var levels []*domain.Level

	for i:=1; i <= cfg.ParkingLevels; i++ {
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

		level,err := domain.NewLevel(i,smallPool,mediumPool,largePool)
		if err != nil {
			return nil,err 
		}

		levels = append(levels, level)
	}

	return &domain.ParkingLot{Levels: levels},nil 
}

func generateSlot(n int) []int {
	slots := make([]int,n)
	for i := 0; i < n; i++ {
		slots[i] = i + 1
	}
	return slots
}