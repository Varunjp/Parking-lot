package domain

type Level struct {
	ID          int
	SmallSlots  *SlotPool
	MediumSlots *SlotPool
	LargeSlots  *SlotPool
}

func (l *Level) GetPool(vType VehicleType) (*SlotPool, bool) {
	switch vType {
	case Small:
		return l.SmallSlots, true
	case Medium:
		return l.MediumSlots, true
	case Large:
		return l.LargeSlots, true
	default:
		return nil, false
	}
}

func NewLevel(id int, small, medium, large *SlotPool) (*Level, error) {
	total := small.Capacity() + medium.Capacity() + large.Capacity()

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