package models

import ("time" 
"github.com/satori/go.uuid")

type MyAgenda struct {
	Base
	UserID    uuid.UUID
	Title     string
	ConferenceID     uuid.UUID
	StartTime time.Time
	EndTime   time.Time

	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt time.Time
}