package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func getReportGeneratorURL() string {
	url := os.Getenv("REPORT_GENERATOR_URL")
	if url == "" {
		url = "http://report-generator.cell-b.svc.cluster.local:8080"
	}
	return url
}
func main() {
	log.Println("Starting analytics-frontend on port 8080...")
	reportURL := getReportGeneratorURL()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("analytics-frontend: Received request, calling report-generator.")
		resp, err := http.Get(reportURL)
		if err != nil {
			http.Error(w, "Failed to connect to report-generator", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()
		var reportData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&reportData); err != nil {
			log.Printf("Error decoding report-generator response: %v", err)
			reportData = map[string]interface{}{"error": "Failed to decode response"}
		}
		hostname, _ := os.Hostname()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":           "Analytics request processed by analytics-frontend",
			"reporter_hostname": hostname,
			"report_data":       reportData,
		})
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
