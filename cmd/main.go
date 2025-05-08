package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"stock-picker/server"
	"stock-picker/stock"
)

// Health status response
type healthStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// healthCheckHandler handles requests to the /healthz endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a health status response
	status := healthStatus{
		Status:    "OK",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the HTTP status code 200 OK
	w.WriteHeader(http.StatusOK)

	// Encode the status struct to JSON and write it to the response
	if err := json.NewEncoder(w).Encode(status); err != nil {
		// If encoding fails, log the error and send an internal server error
		log.Printf("Error encoding health status: %v", err)
		http.Error(w, "Error encoding health status", http.StatusInternalServerError)
	}
}

// readinessCheckHandler handles requests to the /readyz endpoint
func readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	isReady := true

	if isReady {
		status := healthStatus{
			Status:    "READY",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	} else {
		status := healthStatus{
			Status:    "NOT_READY",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(status)
	}
}

func stockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Make sure we get the expected input - ex: localhost:8080/stock/
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "stock" {
		http.Error(w, "Invalid request format. Use localhost:8080/stock/", http.StatusBadRequest)
		return
	}

	// Get stock data from external source
	stockInfo, err := stock.GetStockData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Output the result in json format
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(stockInfo); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// New server instance
	server := server.NewServer(":8080")

	// Health check handler
	server.Handlefunc("/healthz", healthCheckHandler)

	// Readiness check handler
	server.Handlefunc("/readyz", readinessCheckHandler)

	// Specified handler above
	server.Handlefunc("/stock/", stockHandler)

	log.Printf("Starting server...")

	// Start server and listing for requests
	err := server.Start()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Error starting server: ", err)
	}
}
