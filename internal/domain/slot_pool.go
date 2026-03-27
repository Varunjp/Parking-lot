package domain

// Allocate assigns the next available slot from the pool.
//
// Returns:
//   - slot ID (int): the allocated slot
//   - error: ErrNoSlotsAvailable if no slots are free
//
// Behavior:
//   - Uses FIFO strategy (first available slot is allocated)
//   - Moves slot from FreeSlots → Occupied
//
// Time Complexity: O(1)
func (s *SlotPool) Allocate() (int, error) {
	if len(s.FreeSlots) == 0 {
		return 0, ErrNoSlotsAvailable
	}

	// Take the first available slot (FIFO)
	slot := s.FreeSlots[0]

	// Remove slot from free list
	s.FreeSlots = s.FreeSlots[1:]

	// Mark slot as occupied
	s.Occupied[slot] = true

	return slot, nil
}

// Release frees a previously occupied slot.
//
// Parameters:
//   - slot (int): the slot ID to release
//
// Behavior:
//   - Removes slot from Occupied
//   - Adds slot back to FreeSlots queue
//   - If slot is not occupied, operation is ignored (idempotent)
//
// Time Complexity: O(1)
func (s *SlotPool) Release(slot int) {
	if !s.Occupied[slot] {
		// Slot is either already free or invalid → ignore
		return
	}

	// Remove from occupied set
	delete(s.Occupied, slot)

	// Add back to free slots queue
	s.FreeSlots = append(s.FreeSlots, slot)
}

// AvailableCount returns the number of currently free slots.
//
// Returns:
//   - int: count of free slots
//
// Time Complexity: O(1)
func (s *SlotPool) AvailableCount() int {
	return len(s.FreeSlots)
}

// Capacity returns the total number of slots managed by the pool.
//
// Includes both free and occupied slots.
//
// Returns:
//   - int: total slot capacity
//
// Time Complexity: O(1)
func (s *SlotPool) Capacity() int {
	return len(s.FreeSlots) + len(s.Occupied)
}