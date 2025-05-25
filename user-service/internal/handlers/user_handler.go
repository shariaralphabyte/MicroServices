package handlers

import (
	"encoding/json"
	"net/http"
	"shared/events"
	"strconv"
	"userservice/internal/db"
	"userservice/internal/models"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

type UserHandler struct {
	db   *db.PostgresDB
	nats *nats.Conn
}

func NewUserHandler(db *db.PostgresDB, nats *nats.Conn) *UserHandler {
	return &UserHandler{db: db, nats: nats}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish user created event
	event := events.UserEvent{
		Type: "user_created",
		Payload: events.UserPayload{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			UpdatedAt: user.UpdatedAt,
		},
	}
	eventJSON, _ := json.Marshal(event)
	h.nats.Publish("user.events", eventJSON)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = id

	if err := h.db.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish user updated event
	event := events.UserEvent{
		Type: "user_updated",
		Payload: events.UserPayload{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			UpdatedAt: user.UpdatedAt,
		},
	}
	eventJSON, _ := json.Marshal(event)
	h.nats.Publish("user.events", eventJSON)

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
