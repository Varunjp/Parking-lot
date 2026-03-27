package tests

import "parking-lot/internal/domain"

func createLevel(id int, slots int) *domain.Level {
	return &domain.Level{
		ID: id,
		SmallSlots: &domain.SlotPool{
			FreeSlots: generateSlots(slots),
			Occupied:  make(map[int]bool),
		},
		MediumSlots: &domain.SlotPool{
			FreeSlots: generateSlots(slots),
			Occupied:  make(map[int]bool),
		},
		LargeSlots: &domain.SlotPool{
			FreeSlots: generateSlots(slots),
			Occupied:  make(map[int]bool),
		},
	}
}

func generateSlots(n int) []int {
	slots := make([]int, n)
	for i := 0; i < n; i++ {
		slots[i] = i + 1
	}
	return slots
}