package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Config structure
type Config struct {
	SonnenToken   string `json:"sonnen_token"`
	SonnenHost    string `json:"sonnen_host"`
	MypvToken     string `json:"mypv_token"`
	MypvSerial    string `json:"mypv_serial"`
	Interval      string `json:"interval"`
	ReservedWatts int    `json:"reserved"`
	MaxTemp       int    `json:"max_temp"`
	DatabasePath  string `json:"database_path"`
}

var config Config

var ConfigLock sync.Mutex

func GetConfig() Config {
	return config
}

// Load configuration from file
func LoadConfig() error {
	ConfigLock.Lock()
	defer ConfigLock.Unlock()
	file, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &config)
}

// Save configuration to file
func SaveConfig(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request...")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fmt.Println("Received JSON:", string(bodyBytes))

	ConfigLock.Lock()
	defer ConfigLock.Unlock()

	// Use a map to temporarily hold JSON data
	var tempConfig map[string]string
	if err := json.Unmarshal(bodyBytes, &tempConfig); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Manually convert string values to integers where necessary
	config.SonnenToken = tempConfig["sonnen_token"]
	config.SonnenHost = tempConfig["sonnen_host"]
	config.MypvToken = tempConfig["mypv_token"]
	config.MypvSerial = tempConfig["mypv_serial"]
	config.Interval = tempConfig["interval"]

	if reserved, err := strconv.Atoi(tempConfig["reserved"]); err == nil {
		config.ReservedWatts = reserved
	} else {
		http.Error(w, "Invalid reserved value", http.StatusBadRequest)
		return
	}

	if maxTemp, err := strconv.Atoi(tempConfig["max_temp"]); err == nil {
		config.MaxTemp = maxTemp
	} else {
		http.Error(w, "Invalid max_temp value", http.StatusBadRequest)
		return
	}

	// Save the corrected config
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}
	_ = os.WriteFile("config.json", file, 0644)

	w.WriteHeader(http.StatusOK)
}

func ReturnConfig(w http.ResponseWriter, r *http.Request) {
	err := LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	config := GetConfig()
	jsonData, err := json.Marshal(config)
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println("Error writing response:", err)
	}

}

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ReturnConfig(w, r)
	} else if r.Method == http.MethodPost {
		SaveConfig(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}
