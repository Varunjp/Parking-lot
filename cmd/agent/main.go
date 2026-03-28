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

type ExitRequest struct {
	VehicleID string `json:"vehicle_id"`
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

		cmd,data,err := parseInput(input)
		if err != nil {
			fmt.Println("invalid input:",err)
			continue
		}
		switch cmd {

		case "PARK":
			req := data.(ParkRequest)
			if err := sendParkRequest(req,port); err != nil {
				fmt.Println("Request failed:",err)
			}
		
		case "EXIT":
			req := data.(ExitRequest)
			if err := sendExitRequest(req,port); err != nil {
				fmt.Println("Request failed:",err)
			}
		
		default:
			fmt.Println("Unknown command")
		}
	}
}

// sendParkRequest sends the allocation request to the parking service.
// It handles HTTP communication, response decoding, and error reporting.
func sendParkRequest(req ParkRequest,port string) error {

	body,err:= json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %w",err)
	}

	url := fmt.Sprintf("http://localhost:%s/park",port)
	
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

func sendExitRequest(req ExitRequest,port string) error {

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %w", err)
	}

	url := fmt.Sprintf("http://localhost:%s/exit", port)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if status, ok := result["status"].(string); ok && status == "error" {
		return fmt.Errorf("%v", result["message"])
	}

	fmt.Println("Vehicle exited successfully")
	return nil
}

// printInstructions displays CLI usage guidance to the user.
// Keeping this separate improves readability of main().
func printInstructions() {
	fmt.Println("Parking Agent Started")
	fmt.Println("\nCommands:")

	fmt.Println("PARK vehicleID vehicleType customerType")
	fmt.Println("  VehicleType  -> SMALL | MEDIUM | LARGE")
	fmt.Println("  CustomerType -> EMERGENCY | VIP | REGULAR")
	fmt.Println("  Example: PARK KL01 MEDIUM VIP")

	fmt.Println("\nEXIT vehicleID")
	fmt.Println("  Example: EXIT KL01")
}

// parseInput converts raw CLI input into a structured ParkRequest.
// It ensures input format correctness and normalizes values.
func parseInput(input string) (string,interface{}, error) {
	parts := strings.Fields(input) // safer than Split (handles extra spaces)

	if len(parts) == 0 {
		return "",nil, fmt.Errorf("empty input")
	}

	cmd := strings.ToUpper(parts[0])

	switch cmd {
	case "PARK":
		if len(parts) != 4 {
			return "",nil,fmt.Errorf("usage: PARK vehicleID vehicleType customerType")
		}

		return "PARK", ParkRequest{
			VehicleID: parts[1],
			VehicleType: strings.ToUpper(parts[2]),
			CustomerType: strings.ToUpper(parts[3]),
		},nil 
		
	case "EXIT":
		if len(parts) != 2 {
			return "",nil,fmt.Errorf("usage: EXIT vehicleID")
		}

		return "EXIT",ExitRequest{
			VehicleID: parts[1],
		},nil 
	}

	return "",nil,fmt.Errorf("invalid command")
}