package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getUserServiceURL() string {
	url := os.Getenv("USER_SVC_URL")
	if url == "" {
		url = "http://user-svc.cell-a.svc.cluster.local:8080"
	}
	return url
}
func main() {
	log.Println("Starting report-generator on port 8080...")
	userURL := getUserServiceURL()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("report-generator: Received request, fetching data from user-svc in cell-a.")
		resp, err := http.Get(userURL)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to user-svc in cell-a: %v", err)
			http.Error(w, msg, http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()
		var userData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
			log.Printf("Error decoding user-svc response: %v", err)
			userData = map[string]interface{}{"error": "Failed to decode response"}
		}
		hostname, _ := os.Hostname()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"report_status":      "complete",
			"data_source_cell":   "cell-a",
			"generator_hostname": hostname,
			"raw_user_data":      userData,
		})
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
