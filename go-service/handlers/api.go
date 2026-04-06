package handlers

import (
	"campus-connect/go-service/webhook"
	"encoding/json"
	"net/http"
)

func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	webhook.StatsMutex.Lock()
	defer webhook.StatsMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhook.AnalyticsStats)
}

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	webhook.NotifMutex.Lock()
	defer webhook.NotifMutex.Unlock()

	response := []map[string]interface{}{}
	for i, msg := range webhook.Notifications {
		response = append(response, map[string]interface{}{
			"id":      i + 1,
			"message": msg,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
