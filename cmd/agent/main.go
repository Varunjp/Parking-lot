package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"parking-lot/internal/config"
	"strings"
)

// ParkRequest represents the payload sent to the parking service.
// It defines the minimal required fields for slot allocation.
type ParkRequest struct {
	VehicleID    string `json:"vehicle_id"`
	VehicleType  string `json:"vehicle_type"`
	CustomerType string `json:"customer_type"`
}

func main() {
	cfg := config.Load()
	port := cfg.HttpPort

	reader := bufio.NewReader(os.Stdin)

	printInstructions()

	for {
		fmt.Print("> ")

		input,err:= reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read input:",err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue // Ignore empty input
		}

		req,err := parseInput(input)
		if err != nil {
			fmt.Println("invalid input:",err)
			continue
		}

		if err := sendRequest(req,port); err != nil {
			fmt.Println("Request failed:",err)
		}
	}
}

// sendParkRequest sends the allocation request to the parking service.
// It handles HTTP communication, response decoding, and error reporting.
func sendRequest(req ParkRequest,port string) error {

	body,err:= json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %w",err)
	}

	url := fmt.Sprintf("https://localhost:%s/park",port)
	
	resp,err := http.Post(url,"application/json",bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// Decode response into a generic map (can be replaced with a typed struct later)
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w",err)
	}

	// Handle application-level error response
	if status,ok := result["status"].(string); ok && status == "error" {
		return fmt.Errorf("%v",result["message"])
	}

	fmt.Printf("Allocated: Level %v Slot %v\n",result["level"],result["slot"])
	return nil 
}

// printInstructions displays CLI usage guidance to the user.
// Keeping this separate improves readability of main().
func printInstructions() {
	fmt.Println("Parking Agent Started")
	fmt.Println("Enter: vehicleID vehicleType customerType")
	fmt.Println("VehicleType  -> SMALL | MEDIUM | LARGE")
	fmt.Println("CustomerType -> EMERGENCY | VIP | REGULAR")
	fmt.Println("Example: KL01 MEDIUM VIP")
}

// parseInput converts raw CLI input into a structured ParkRequest.
// It ensures input format correctness and normalizes values.
func parseInput(input string) (ParkRequest, error) {
	parts := strings.Fields(input) // safer than Split (handles extra spaces)

	if len(parts) != 3 {
		return ParkRequest{}, fmt.Errorf("expected 3 fields: vehicleID vehicleType customerType")
	}

	return ParkRequest{
		VehicleID:    parts[0],
		VehicleType:  strings.ToUpper(parts[1]),
		CustomerType: strings.ToUpper(parts[2]),
	}, nil
}