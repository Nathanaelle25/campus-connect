package webhook

import (
	"campus-connect/go-service/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	AnalyticsStats = map[string]int{"total_events": 0, "total_webhooks_received": 0}
	StatsMutex     sync.Mutex
	Notifications  []string
	NotifMutex     sync.Mutex
)

// ProcessWebhook receives the webhook and fires multiple goroutines
func ProcessWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload models.WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Respond immediately that it was accepted
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "processing"}`))

	// Bonus B: Concurrent operations using Goroutines and WaitGroup
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		// Mock DB Write
		time.Sleep(10 * time.Millisecond)
		NotifMutex.Lock()
		Notifications = append(Notifications, fmt.Sprintf("Event '%s' processed at %s", payload.Title, payload.Timestamp))
		NotifMutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		// Log to file (simulated)
		time.Sleep(5 * time.Millisecond)
	}()

	go func() {
		defer wg.Done()
		// Update Analytics
		StatsMutex.Lock()
		AnalyticsStats["total_events"]++
		AnalyticsStats["total_webhooks_received"]++
		StatsMutex.Unlock()
	}()

	wg.Wait()
}
