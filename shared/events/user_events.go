package events

type UserEvent struct {
	Type    string      `json:"type"`
	Payload UserPayload `json:"payload"`
}

type UserPayload struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	UpdatedAt string `json:"updated_at"`
}
