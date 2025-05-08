package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"stock-picker/server"
	"stock-picker/stock"
)

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

	// Specified handler above
	server.Handlefunc("/stock/", stockHandler)

	log.Printf("Starting server...")

	// Start server and listing for requests
	err := server.Start()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Error starting server: ", err)
	}
}
