package domain

// VehicleType represents the category of a vehicle based on its size.
// It is used to determine suitable parking slot allocation.
type VehicleType string

// CustomerType represents the priority category of a customer.
// This is used in allocation strategies (e.g., emergency > VIP > regular).
type CustomerType string

const (
	// Vehicle Types

	// Small represents compact vehicles like bikes or small cars.
	// Typically allocated to small-sized parking slots.
	Small  VehicleType = "SMALL"

	// Medium represents standard vehicles like sedans.
	// Can fit in medium or larger slots.
	Medium VehicleType = "MEDIUM"

	// Large represents bigger vehicles like SUVs or trucks.
	// Requires larger parking slots.
	Large  VehicleType = "LARGE"

	// Customer Types

	// Regular represents normal customers with standard priority.
	Regular   CustomerType = "REGULAR"

	// VIP represents high-priority customers.
	// May get preferential slot allocation.
	VIP       CustomerType = "VIP"

	// Emergency represents highest-priority vehicles such as ambulances,
	// fire trucks, or police vehicles. These should be allocated slots first.
	Emergency CustomerType = "EMERGENCY"
)
