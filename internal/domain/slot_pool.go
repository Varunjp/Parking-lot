package domain

type SlotPool struct {
	FreeSlots []int
	Occupied  map[int]bool
}

func (s *SlotPool) Allocate() (int, error) {
	if len(s.FreeSlots) == 0 {
		return 0, ErrNoSlotsAvailable
	}

	slot := s.FreeSlots[0]
	s.FreeSlots = s.FreeSlots[1:]
	s.Occupied[slot] = true
	return slot, nil
}

func (s *SlotPool) Release(slot int) {
	if s.Occupied[slot] {
		delete(s.Occupied, slot)
		s.FreeSlots = append(s.FreeSlots, slot)
	}
}

func (s *SlotPool) AvailableCount() int {
	return len(s.FreeSlots)
}

func (s *SlotPool) Capacity() int {
	return len(s.FreeSlots) + len(s.Occupied)
}