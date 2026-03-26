package domain

import "errors"

var (
	ErrParkingFull 	= errors.New("parking full")
	ErrReEntry  	= errors.New("re-entry restricted")
	ErrInvalidType 	= errors.New("invalid vehicle type")
)