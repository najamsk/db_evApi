package models

import "github.com/satori/go.uuid"

type StartupLevel struct {
	Base
	Name             string
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	SortOrder int
	Startups         []Startup

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}