package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// This service calls the report-generator in the same cell.
func getReportGeneratorURL() string {
	url := os.Getenv("REPORT_GENERATOR_URL")
	if url == "" {
		// This URL points to the report-generator service within cell-b
		url = "http://report-generator.cell-b.svc.cluster.local:8080"
	}
	return url
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("analytics-frontend: Received request, calling report-generator.")

	reportURL := getReportGeneratorURL()
	resp, err := http.Get(reportURL)
	if err != nil {
		log.Printf("Error calling report-generator: %v", err)
		http.Error(w, "Failed to connect to report-generator", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var reportData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&reportData)

	log.Println("analytics-frontend: Received report, preparing final response.")

	hostname, _ := os.Hostname()
	finalResponse := map[string]interface{}{
		"message":           "Analytics request processed by analytics-frontend",
		"reporter_hostname": hostname,
		"report_data":       reportData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalResponse)
	log.Println("analytics-frontend: Response sent.")
}

func main() {
	log.Println("Starting analytics-frontend on port 8080...")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start analytics-frontend: %v", err)
	}
}
