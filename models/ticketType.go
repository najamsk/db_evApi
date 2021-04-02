package models

import ("github.com/satori/go.uuid"
"time")

type TicketType struct {
	Base
	Title           string
	IsActive        bool
	ClientID        uuid.UUID
	ConferenceID    uuid.UUID
	Amount          float64
	AmmountCurrency string
	Tickets         []Ticket
	Description		string
	ValidFrom       time.Time
	ValidTo         time.Time
	MaxQuantityBook	int
	
	// CreatedAt       time.Time
	// UpdatedAt       time.Time
	// DeletedAt       time.Time
}