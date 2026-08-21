package events

import "time"

type EventType string

const (
	BookCreated EventType = "book.created" // implemented
	BookUpdated EventType = "book.updated" // implemented
	BookDeleted EventType = "book.deleted" // implemented

	UserRegistered EventType = "user.registered" // implemented
	UserUpdated    EventType = "user.updated"
	UserLoggedIn   EventType = "user.loggedIn" // implemented
	UserDeleted    EventType = "user.deleted"
)

type Envelope struct {
	EventID    string    `json:"event_id"`
	EventType  EventType `json:"event_type"`
	OccurredAt time.Time `json:"occured_at"`
	EntityID   int64     `json:"entity_id"`
	Payload    any       `json:"payload"`
}

type BookCreatedPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorID    int64  `json:"author_id"`
}

type BookUpdatedPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type BookDeletedPayload struct {
}

type UserRegisteredPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IP        string `json:"ip,omitempty"`
}

type UserUpdatedPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserLoggedInPayload struct {
	IP string `json:"ip,omitempty"`
}

type UserDeletedPayload struct {
}
