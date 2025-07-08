package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Message  string `json:"message"`
	UserData string `json:"user_data"`
	AuthData string `json:"auth_data"`
	Hostname string `json:"hostname"`
}

// Get the user-svc URL from an environment variable, with a default
func getUserServiceURL() string {
	url := os.Getenv("USER_SVC_URL")
	if url == "" {
		url = "http://user-svc.cell-a.svc.cluster.local:8080" // Default for testing
	}
	return url
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("frontend-svc: Received request, calling user-svc.")

	// Call the downstream user-svc
	userURL := getUserServiceURL()
	resp, err := http.Get(userURL)
	if err != nil {
		log.Printf("Error calling user-svc: %v", err)
		http.Error(w, "Failed to connect to user-svc", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		log.Printf("Error decoding user-svc response: %v", err)
		http.Error(w, "Failed to decode user-svc response", http.StatusInternalServerError)
		return
	}

	log.Println("frontend-svc: Received response from user-svc, preparing final response.")

	hostname, _ := os.Hostname()
	
	// Safe type assertions with defaults
	userData := "N/A"
	if msg, ok := userResponse["message"].(string); ok {
		userData = msg
	}
	
	authData := "N/A"
	if auth, ok := userResponse["auth_data"].(string); ok {
		authData = auth
	}
	
	finalResponse := Response{
		Message:  "Request processed by frontend-svc",
		UserData: userData,
		AuthData: authData,
		Hostname: hostname,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResponse)
	log.Println("frontend-svc: Response sent.")
}

func main() {
	log.Println("Starting frontend-svc on port 8080...")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start frontend-svc: %v", err)
	}
}
