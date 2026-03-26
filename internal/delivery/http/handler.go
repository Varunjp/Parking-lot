package httpClient

import (
	"encoding/json"
	"net/http"
	"parking-lot/internal/domain"
	"parking-lot/internal/usecase"
)

type Handler struct {
	dispatcher *usecase.Dispatcher
}

func NewHandler(d *usecase.Dispatcher) *Handler {
	return &Handler{dispatcher: d}
}

type Request struct {
	VehicleID 		string `json:"vehicle_id"`
	VehicleType    	string `json:"vehicle_type"`
	CustomerType	string `json:"customer_type"`
}

type Response struct {
	Status  	string `json:"status"`
	Message     string `json:"message,omitempty"`
	Level       int    `json:"level,omitempty"`
	Slot        int    `json:"slot,omitempty"`
}

func (h *Handler) Park(w http.ResponseWriter,r *http.Request) {

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "invalid request",
		})
		return 
	}

	v := domain.Vehicle {
		ID: 			req.VehicleID,
		Type: 			domain.VehicleType(req.VehicleType),
		CustomerType: 	domain.CustomerType(req.CustomerType),
	}

	result := h.dispatcher.AddRequest(v)

	if result.Err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: result.Err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Status: "success",
		Level:  result.Level,
		Slot:   result.Slot,
	})
}

func writeJSON(w http.ResponseWriter,code int,res Response) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}