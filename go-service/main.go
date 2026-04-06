package main

import (
	"campus-connect/go-service/handlers"
	"campus-connect/go-service/middleware"
	"campus-connect/go-service/webhook"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Standard http mux
	mux := http.NewServeMux()

	// Endpoints that require both Rate Limiting and API Key
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/analytics", handlers.AnalyticsHandler)
	protectedMux.HandleFunc("/notifications", handlers.NotificationsHandler)

	// Wrap with middlewares
	protectedHandler := middleware.RateLimitMiddleware(middleware.APIKeyAuthMiddleware(protectedMux))

	// Map routes
	// Root routes that will go mapped
	mux.Handle("/analytics", protectedHandler)
	mux.Handle("/notifications", protectedHandler)
	
	// Webhook endpoint (Waitgroup goroutine processing internal)
	mux.HandleFunc("/webhook", webhook.ProcessWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Go Service starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
