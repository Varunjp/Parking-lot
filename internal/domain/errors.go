package domain

import "errors"

// -----------------------------------------------------------------------------
// Domain Errors
// -----------------------------------------------------------------------------
//
// This file defines all domain-level errors used across the parking system.
// These errors represent business rule violations and invalid operations.
//
// Guidelines:
// - Use these as sentinel errors for comparison (errors.Is)
// - Keep messages user-friendly but not UI-specific
// - Add context at service/handler layer using fmt.Errorf("%w", err)
// -----------------------------------------------------------------------------

// Sentinel errors for common domain failures.
var (
	// ErrParkingFull indicates that the parking lot has reached maximum capacity.
	ErrParkingFull 	= errors.New("parking full")

	// ErrReEntry indicates that a vehicle is attempting to re-enter
	// within a restricted time window for the same level.
	ErrReEntry  	= errors.New("no valid level available due to re-entry restriction")

	// ErrNoSlotsAvailable indicates that no suitable parking slots are available
	// for the given vehicle type or constraints.
	ErrNoSlotsAvailable = errors.New("no slots available")

	// ErrInvalidType indicates that the provided vehicle type is not supported.
	ErrInvalidType 	= errors.New("invalid vehicle type")

	// ErrInvalidPriority indicates that the vehicle priority is not recognized.
	ErrInvalidPriority = errors.New("invalid priority")

	// ErrInvalidCapacity indicates that the parking level configuration is invalid.
	ErrInvalidCapacity = errors.New("parking level capacity invalid")
)

// -----------------------------------------------------------------------------
// Typed Errors (for structured validation failures)
// -----------------------------------------------------------------------------

// DomainError represents a custom error type for validation and request-level issues.
// It allows differentiation from system/internal errors.
type domainError string 

var (
	// ErrInvalidRequest indicates that the incoming request is malformed
	// or does not meet required validation criteria.
	ErrInvalidRequest domainError = "invalid request"
)