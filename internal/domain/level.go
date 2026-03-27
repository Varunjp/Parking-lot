package domain

// GetPool returns the SlotPool corresponding to the given VehicleType.
//
// It maps a vehicle type (Small, Medium, Large) to its respective slot pool
// within the level. The boolean return value indicates whether a valid pool
// was found for the given type.
//
// Returns:
//   - *SlotPool: the matching slot pool for the vehicle type
//   - bool: false if the vehicle type is unsupported
func (l *Level) GetPool(vType VehicleType) (*SlotPool, bool) {
	switch vType {
	case Small:
		return l.SmallSlots, true
	case Medium:
		return l.MediumSlots, true
	case Large:
		return l.LargeSlots, true
	default:
		// Unknown vehicle type — caller should handle this case explicitly
		return nil, false
	}
}

// NewLevel creates and initializes a new Level with the provided slot pools.
//
// It validates the total capacity of the level by summing the capacities
// of all slot pools (Small, Medium, Large). A valid level must have a total
// capacity within the allowed range [10, 100].
//
// Parameters:
//   - id: unique identifier for the level
//   - small: slot pool for small vehicles
//   - medium: slot pool for medium vehicles
//   - large: slot pool for large vehicles
//
// Returns:
//   - *Level: initialized level instance if validation passes
//   - error: ErrInvalidCapacity if total slots are outside allowed limits
//
// Notes:
//   - Assumes SlotPool is non-nil and Capacity() is safe to call
//   - Caller is responsible for ensuring slot pool correctness
func NewLevel(id int, small, medium, large *SlotPool) (*Level, error) {
	// Calculate total capacity across all slot types
	total := small.Capacity() + medium.Capacity() + large.Capacity()

	// Enforce business constraint on level capacity
	if total < 10 || total > 100 {
		return nil, ErrInvalidCapacity
	}

	return &Level{
		ID:          id,
		SmallSlots:  small,
		MediumSlots: medium,
		LargeSlots:  large,
	}, nil
}