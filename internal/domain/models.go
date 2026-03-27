package domain

// Vehicle represents a vehicle entering the parking system.
// It contains identity and classification details used for allocation decisions.
type Vehicle struct {
	// ID is a unique identifier for the vehicle (e.g., license plate).
	ID           string

	// Type defines the physical size category of the vehicle
	// (e.g., Small, Medium, Large) and determines slot eligibility.
	Type         VehicleType

	// CustomerType indicates priority level (e.g., Regular, VIP, Emergency).
	// This can influence allocation strategy and slot preference.
	CustomerType CustomerType
}

// SlotPool manages a group of slots of the same type (Small/Medium/Large).
// It is optimized for O(1) allocation and deallocation.
type SlotPool struct {
	// FreeSlots holds available slot IDs.
	// Typically implemented as a stack or queue for efficient access.
	FreeSlots []int

	// Occupied tracks currently occupied slot IDs.
	// The map provides O(1) lookup for validation and release operations.
	Occupied  map[int]bool
}

// ParkingLot represents the entire parking facility.
// It consists of multiple levels.
type ParkingLot struct {
	// Levels contains all parking levels in the system.
	Levels []*Level
}

// Level represents a single floor in the parking lot.
// Each level has separate slot pools for different vehicle sizes.
type Level struct {
	// ID uniquely identifies the level (e.g., Floor number).
	ID          int
	
	// SmallSlots manages slots for small vehicles (e.g., bikes).
	SmallSlots  *SlotPool

	// MediumSlots manages slots for medium vehicles (e.g., cars).
	MediumSlots *SlotPool

	// LargeSlots manages slots for large vehicles (e.g., trucks).
	LargeSlots  *SlotPool
}
