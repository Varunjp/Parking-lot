package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values loaded from environment variables.
// It provides default values when environment variables are not set.
type Config struct {
	ParkingLevels  			int		// Total number of parking levels in the system
	SmallSlotsPerLevel		int 	// Number of small vehicle slots per level
	MediumSlotsPerLevel     int 	// Number of medium vehicle slots per level
	LargeSlotsPerLevel      int 	// Number of large vehicle slots per level
	ReEntrySeconds 			int64	// Time restriction (in seconds) for vehicle re-entry
	HttpPort       			string 	// Port on which HTTP server will run
}

// Load initializes and returns the application configuration.
//
// It attempts to load environment variables from a `.env` file.
// If the file is not found, it falls back to system environment variables
// and default values defined in the helper functions.
func Load() *Config {

	// Attempt to load .env file (non-critical)
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using defaults")
	}

	return &Config{
		ParkingLevels: getInt("PARKING_LEVELS",2),
		SmallSlotsPerLevel: getInt("SMALL_SLOTS_PER_LEVEL",5),
		MediumSlotsPerLevel: getInt("MEDIUM_SLOTS_PER_LEVEL",3),
		LargeSlotsPerLevel: getInt("LARGE_SLOTS_PER_LEVEL",2),
		ReEntrySeconds: int64(getInt("REENTRY_SECONDS",3600)),
		HttpPort: getVal("HTTP_PORT","8080"),
	}
}

// getInt retrieves an integer value from environment variables.
//
// If the environment variable is not set or cannot be parsed as an integer,
// it returns the provided default value.
//
// Example:
//   PARKING_LEVELS=3 → returns 3
//   PARKING_LEVELS=abc → returns default value
func getInt(key string,defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val,err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val 
}

// getVal retrieves a string value from environment variables.
//
// If the environment variable is not set, it returns the provided default value.
//
// Example:
//   HTTP_PORT=9090 → returns "9090"
//   (not set) → returns default value
func getVal(key string,defaultVal string)string {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	return valStr
}