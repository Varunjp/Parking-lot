package httpClient

import (
	"encoding/json"
	"errors"
	"net/http"
	"parking-lot/internal/domain"
	"parking-lot/internal/usecase"
)

// Handler handles HTTP requests related to parking operations.
// It acts as a bridge between HTTP layer and business logic (usecase layer).
type Handler struct {
	dispatcher *usecase.Dispatcher
}

// NewHandler initializes and returns a new Handler instance.
func NewHandler(d *usecase.Dispatcher) *Handler {
	return &Handler{dispatcher: d}
}

// ParkRequest represents the incoming JSON payload for parking a vehicle.
type ParkRequest struct {
	VehicleID 		string `json:"vehicle_id"`
	VehicleType    	string `json:"vehicle_type"`
	CustomerType	string `json:"customer_type"`
}

type ExitRequest struct {
	VehicleID string `json:"vehicle_id"`
}

// ParkResponse represents the API response returned to the client.
type ParkResponse struct {
	Status  	string `json:"status"`				// "success" or "error"
	Message     string `json:"message,omitempty"`	// Error message (if any)
	Level       int    `json:"level,omitempty"`		// Allocated level
	Slot        int    `json:"slot,omitempty"`		// Allocated slot
}

// Park handles vehicle parking requests.
// Flow:
// 1. Decode request body
// 2. Validate input
// 3. Convert to domain model
// 4. Call dispatcher (business logic)
// 5. Return structured JSON response
func (h *Handler) Park(w http.ResponseWriter,r *http.Request) {
	var req ParkRequest

	// Step 1: Decode JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ParkResponse{
			Status:  "error",
			Message: "invalid request",
		})
		return 
	}

	// Step 2: Basic validation (prevents bad data entering business layer)
	if err := validateRequest(req); err != nil {
		writeJSON(w, http.StatusBadRequest, ParkResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Step 3: Map request to domain model
	vehicle := domain.Vehicle {
		ID: 			req.VehicleID,
		Type: 			domain.VehicleType(req.VehicleType),
		CustomerType: 	domain.CustomerType(req.CustomerType),
	}

	// Step 4: Call business logic (dispatcher)
	result := h.dispatcher.Park(vehicle)

	// Step 5: Handle business errors
	if result.Err != nil {
		writeJSON(w, http.StatusBadRequest, ParkResponse{
			Status:  "error",
			Message: result.Err.Error(),
		})
		return
	}

	// Step 6: Success response
	writeJSON(w, http.StatusOK, ParkResponse{
		Status: "success",
		Level:  result.Level,
		Slot:   result.Slot,
	})
}

func (h *Handler) Exit(w http.ResponseWriter,r *http.Request) {
	var req ExitRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w,http.StatusBadRequest,ParkResponse{
			Status: "error",
			Message: "invalid request",
		})
		return 
	}

	if req.VehicleID == "" {
		writeJSON(w,http.StatusBadRequest,ParkResponse{
			Status: "error",
			Message: "vehicle_id is requried",
		})
		return 
	}

	result := h.dispatcher.Exit(req.VehicleID)

	if result.Err != nil {
		writeJSON(w,http.StatusBadRequest,ParkResponse{
			Status: "error",
			Message: result.Err.Error(),
		})
		return 
	}
	msg := "vehicle "+req.VehicleID+" exited successfully"
	writeJSON(w,http.StatusOK,ParkResponse{
		Status: "success",
		Message: msg,
	})
}

// validateRequest performs basic validation on incoming request data.
// Keeps handler clean and separates validation logic.
func validateRequest(req ParkRequest) error {
	if req.VehicleID == "" {
		return errors.New("vehicle_id is required")
	}
	if req.VehicleType == "" {
		return errors.New("vehicle_type is required")
	}
	if req.CustomerType == "" {
		return errors.New("customer_type is required")
	}
	return nil
}

// writeJSON is a helper function to send JSON responses.
// Centralizing response writing ensures consistent headers and encoding.
func writeJSON(w http.ResponseWriter,code int,res ParkResponse) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)

	// Ignoring encode error here since writing to ResponseWriter rarely fails,
	// but in production you may log this.
	json.NewEncoder(w).Encode(res)
}