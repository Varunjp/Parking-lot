package domain

import "errors"

type domainError string 

var (
	ErrParkingFull 	= errors.New("parking full")
	ErrReEntry  	= errors.New("no valid level available due to re-entry restriction")
	ErrNoSlotsAvailable = errors.New("no slots available")
	ErrInvalidType 	= errors.New("invalid vehicle type")
	ErrInvalidPriority = errors.New("invalid priority")
	ErrInvalidCapacity = errors.New("parking level capacity invalid")
)

var (
	ErrInvalidRequest domainError = "invalid request"
)