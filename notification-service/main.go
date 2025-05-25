package main

import (
	"log"
	"net/http"
	"notificationservice/internal/db"
	"notificationservice/internal/handlers"
	"os"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// Get environment variables
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	natsURL := getEnv("NATS_URL", nats.DefaultURL)

	// Initialize MongoDB connection
	mongodb, err := db.NewMongoDB(mongoURI, "notifications_db", "user_notifications")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Initialize handler
	notificationHandler := handlers.NewNotificationHandler(mongodb, nc)

	// Subscribe to user events
	nc.Subscribe("user.events", notificationHandler.HandleUserEvent)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/notifications", notificationHandler.GetAllNotifications).Methods("GET")
	r.HandleFunc("/notifications/user/{id}", notificationHandler.GetUserNotification).Methods("GET")

	// Start server
	log.Printf("Notification service starting on :8081")
	log.Printf("Registered routes: /notifications (GET), /notifications/user/{id} (GET)")
	log.Fatal(http.ListenAndServe(":8081", r))
}
