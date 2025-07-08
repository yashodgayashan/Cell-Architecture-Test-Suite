package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Message  string `json:"message"`
	Hostname string `json:"hostname"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("auth-svc: Received request, preparing response.")

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %v", err)
		hostname = "unknown"
	}

	response := Response{
		Message:  "Authenticated successfully by auth-svc",
		Hostname: hostname,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Println("auth-svc: Response sent.")
}

func main() {
	log.Println("Starting auth-svc on port 8080...")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start auth-svc: %v", err)
	}
}
