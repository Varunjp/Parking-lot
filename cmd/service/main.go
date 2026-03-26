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

	parkingRepo := &memory.ParkingRepo{
		Levels: initLevels(cfg),
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

func initLevels(cfg *config.Config) []*domain.Level {

	var levels []*domain.Level

	for i:=1; i <= cfg.ParkingLevels; i++ {
		levels = append(levels, &domain.Level{
			ID: i,
			SmallSlots: &domain.SlotPool{
				FreeSlots: generateSlot(cfg.SlotsPerLevel),
				Occupied: make(map[int]bool),
			},
			MediumSlots: &domain.SlotPool{
				FreeSlots: generateSlot(cfg.SlotsPerLevel),
				Occupied: make(map[int]bool),
			},
			LargeSlots: &domain.SlotPool{
				FreeSlots: generateSlot(cfg.SlotsPerLevel),
				Occupied: make(map[int]bool),
			},
		})
	}
	return levels
}

func generateSlot(n int) []int {
	slots := make([]int,n)
	for i := 0; i < n; i++ {
		slots[i] = i + 1
	}
	return slots
}