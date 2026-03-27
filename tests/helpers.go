package tests

import "parking-lot/internal/domain"

// createTestLevel constructs a fully initialized parking level for testing.
//
// It creates a Level with independent slot pools for Small, Medium, and Large vehicles.
// Each slot pool is initialized with:
//   - Sequential free slot IDs starting from 1 up to `slots`
//   - An empty occupied map to track allocated slots
//
// This helper ensures consistency across tests and avoids repetitive setup logic.
//
// Parameters:
//   - id: unique identifier for the level
//   - slots: number of slots per vehicle type (Small, Medium, Large)
//
// Returns:
//   - *domain.Level: a ready-to-use level instance for testing
func createLevel(id int, slots int) *domain.Level {
	return &domain.Level{
		ID: id,
		SmallSlots: newSlotPool(slots),
		MediumSlots: newSlotPool(slots),
		LargeSlots: newSlotPool(slots),
	}
}

// newSlotPool initializes a SlotPool with a predefined number of free slots.
//
// It generates slot IDs in sequential order (1..n) and prepares
// an empty Occupied map for tracking allocations during tests.
//
// This abstraction improves readability and avoids duplication
// when creating multiple slot pools.
//
// Parameters:
//   - n: number of slots to initialize
//
// Returns:
//   - *domain.SlotPool: initialized slot pool
func newSlotPool(n int) *domain.SlotPool {
	return &domain.SlotPool{
		FreeSlots: generateSequentialSlots(n),
		Occupied:  make(map[int]bool),
	}
}

// generateSequentialSlots creates a slice of slot IDs from 1 to n.
//
// Example:
//   n = 3 → [1, 2, 3]
//
// This function is used to simulate available parking slots
// in a predictable and deterministic manner, which is important for testing.
//
// Parameters:
//   - n: total number of slots
//
// Returns:
//   - []int: slice containing sequential slot IDs
func generateSequentialSlots(n int) []int {
	slots := make([]int, n)
	for i := 0; i < n; i++ {
		slots[i] = i + 1
	}
	return slots
}