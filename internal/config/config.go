package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ParkingLevels  			int
	SmallSlotsPerLevel		int 
	MediumSlotsPerLevel     int 
	LargeSlotsPerLevel      int 
	ReEntrySeconds 			int64
	HttpPort       			string 
}

func Load() *Config {

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

func getVal(key string,defaultVal string)string {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	return valStr
}