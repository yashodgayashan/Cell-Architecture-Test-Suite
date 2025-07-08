package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// This service calls the user-svc in cell-a to get data.
func getUserServiceURL() string {
	url := os.Getenv("USER_SVC_URL")
	if url == "" {
		// IMPORTANT: This URL points across namespaces to the service in cell-a.
		url = "http://user-svc.cell-a.svc.cluster.local:8080"
	}
	return url
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("report-generator: Received request, fetching data from user-svc in cell-a.")

	// Call the user-svc in cell-a
	userURL := getUserServiceURL()
	resp, err := http.Get(userURL)
	if err != nil {
		log.Printf("Error calling user-svc in cell-a: %v", err)
		http.Error(w, fmt.Sprintf("Failed to connect to user-svc in cell-a: %s", err.Error()), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		log.Printf("Error decoding response from user-svc: %v", err)
		http.Error(w, "Failed to decode response from user-svc", http.StatusInternalServerError)
		return
	}

	log.Println("report-generator: Successfully fetched data. Generating report.")

	// Simple "report generation": count the keys in the response from user-svc
	processedItems := len(userData)

	hostname, _ := os.Hostname()
	reportResponse := map[string]interface{}{
		"report_status":      "complete",
		"data_source_cell":   "cell-a",
		"processed_items":    processedItems,
		"generator_hostname": hostname,
		"raw_user_data":      userData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reportResponse)
	log.Println("report-generator: Report sent.")
}

func main() {
	log.Println("Starting report-generator on port 8080...")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start report-generator: %v", err)
	}
}
