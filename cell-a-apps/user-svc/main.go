package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Message  string `json:"message"`
	AuthData string `json:"auth_data"`
	Hostname string `json:"hostname"`
}

// Get the auth-svc URL from an environment variable, with a default
func getAuthServiceURL() string {
	url := os.Getenv("AUTH_SVC_URL")
	if url == "" {
		url = "http://auth-svc.cell-a.svc.cluster.local:8080" // Default for testing
	}
	return url
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("user-svc: Received request, calling auth-svc.")

	// Call the downstream auth-svc
	authURL := getAuthServiceURL()
	resp, err := http.Get(authURL)
	if err != nil {
		log.Printf("Error calling auth-svc: %v", err)
		http.Error(w, "Failed to connect to auth-svc", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var authResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		log.Printf("Error decoding auth-svc response: %v", err)
		http.Error(w, "Failed to decode auth-svc response", http.StatusInternalServerError)
		return
	}

	log.Println("user-svc: Received response from auth-svc, preparing final response.")

	hostname, _ := os.Hostname()
	
	// Safe type assertion with default
	authData := "N/A"
	if msg, ok := authResponse["message"].(string); ok {
		authData = msg
	}
	
	finalResponse := Response{
		Message:  "User data retrieved by user-svc",
		AuthData: authData,
		Hostname: hostname,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResponse)
	log.Println("user-svc: Response sent.")
}

func main() {
	log.Println("Starting user-svc on port 8080...")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start user-svc: %v", err)
	}
}
