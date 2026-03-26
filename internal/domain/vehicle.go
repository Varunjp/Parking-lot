package domain

type VehicleType string
type CustomerType string

const (
	Small  VehicleType = "SMALL"
	Medium VehicleType = "MEDIUM"
	Large  VehicleType = "LARGE"

	Regular   CustomerType = "REGULAR"
	VIP       CustomerType = "VIP"
	Emergency CustomerType = "EMERGNECY"
)

type Vehicle struct {
	ID           string
	Type         VehicleType
	CustomerType CustomerType
}