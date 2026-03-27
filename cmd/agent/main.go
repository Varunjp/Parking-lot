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

type Request struct {
	VehicleID    string `json:"vehicle_id"`
	VehicleType  string `json:"vehicle_type"`
	CustomerType string `json:"customer_type"`
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	cfg := config.Load()
	port := cfg.HttpPort

	fmt.Println("Agent Started")
	fmt.Println("Enter: vehicleID vehicleType customerType")
	fmt.Println("vehicletype (SMALL,MEDIUM,LARGE)")
	fmt.Println("customerType (EMERGENCY,VIP,REGULAR)")
	fmt.Println("Example: KL01 MEDIUM VIP")

	for {
		fmt.Print("> ")

		input,_ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Split(input," ")

		if len(parts) != 3 {
			fmt.Println("Invalid input format")
			continue
		}

		req := Request {
			VehicleID: parts[0],
			VehicleType: strings.ToUpper(parts[1]),
			CustomerType: strings.ToUpper(parts[2]),
		}

		sendRequest(req,port)
	}
}

func sendRequest(req Request,port string) {

	body,_ := json.Marshal(req)

	resp,err := http.Post(
		"http://localhost:"+port+"/park",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		fmt.Println("Error: ",err)
		return 
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["status"] == "error" {
		fmt.Println("Failed: ",result["message"])
		return 
	}

	fmt.Printf("Allocated: Level %v Slot %v\n",result["level"],result["slot"])
}