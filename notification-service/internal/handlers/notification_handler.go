package handlers

import (
	"encoding/json"
	"net/http"
	"notificationservice/internal/db"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

type NotificationHandler struct {
	db   *db.MongoDB
	nats *nats.Conn
}

func NewNotificationHandler(db *db.MongoDB, nats *nats.Conn) *NotificationHandler {
	return &NotificationHandler{db: db, nats: nats}
}

func (h *NotificationHandler) HandleUserEvent(msg *nats.Msg) {
	var event struct {
		Type    string `json:"type"`
		Payload struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			UpdatedAt string `json:"updated_at"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return
	}

	notification := &db.UserNotification{
		UserID:    event.Payload.ID,
		Name:      event.Payload.Name,
		Email:     event.Payload.Email,
		UpdatedAt: time.Now(),
	}

	h.db.UpsertUser(notification)
}

func (h *NotificationHandler) GetUserNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	notification, err := h.db.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(notification)
}

func (h *NotificationHandler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := h.db.GetAllNotifications()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
