package domain

import "errors"

type domainError string 

var (
	ErrParkingFull 	= errors.New("parking full")
	ErrReEntry  	= errors.New("no valid level available due to re-entry restriction")
	ErrInvalidType 	= errors.New("invalid vehicle type")
	ErrInvalidPriority = errors.New("invalid priority")
)

var (
	ErrInvalidRequest domainError = "invalid request"
)