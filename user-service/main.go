package main

import (
	"log"
	"net/http"
	"os"
	"userservice/internal/db"
	"userservice/internal/handlers"

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
	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPort := getEnv("POSTGRES_PORT", "5432")
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPass := getEnv("POSTGRES_PASSWORD", "password")
	pgDB := getEnv("POSTGRES_DB", "users_db")
	natsURL := getEnv("NATS_URL", nats.DefaultURL)

	// Initialize PostgreSQL connection
	postgres, err := db.NewPostgresDB(pgHost, pgPort, pgUser, pgPass, pgDB)
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
	userHandler := handlers.NewUserHandler(postgres, nc)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")

	// Start server
	log.Printf("User service starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
