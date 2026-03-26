package domain

type SlotPool struct {
	FreeSlots []int
	Occupied  map[int]bool
}

type Level struct {
	ID          int
	SmallSlots  *SlotPool
	MediumSlots *SlotPool
	LargeSlots  *SlotPool
}