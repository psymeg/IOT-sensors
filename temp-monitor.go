package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
        "log/syslog"
)

// TemperatureData represents the JSON payload from the Arduino
type TemperatureData struct {
	Temperature float64 `json:"temperature"`
}

const (
	threshold     = 62.0 // Alert if temperature exceeds this in °C
	listenAddress = ":8080"
)

func main() {
	http.HandleFunc("/sensor", sensorHandler)
	fmt.Printf("Listening for temperature data on %s\n", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func sensorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var data TemperatureData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("[%s] Received temperature: %.2f°C\n",
		time.Now().Format(time.RFC3339), data.Temperature)

	if data.Temperature > threshold {
		triggerAlert(data.Temperature)
	}

	w.WriteHeader(http.StatusOK)
}

func triggerAlert(temp float64) {
	// Connect to syslog
	sysLog, err := syslog.New(syslog.LOG_ALERT|syslog.LOG_USER, "TempMonitor")
	if err != nil {
		fmt.Printf("Error connecting to syslog: %v\n", err)
		return
	}
	defer sysLog.Close()

	msg := fmt.Sprintf("CRITICAL: Temperature exceeded threshold! Current: %.2f°C, Threshold: %.2f°C", temp, threshold)

	// Send to syslog with ALERT priority
	if err := sysLog.Alert(msg); err != nil {
		fmt.Printf("Error writing to syslog: %v\n", err)
	}
	// Print to console
	//fmt.Printf("CRITICAL: Temperature exceeded threshold! Current: %.2f°C, Threshold: %.2f°C\n", temp, threshold)
}
