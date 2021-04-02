package models

import(
	 "github.com/satori/go.uuid"
)

type TicketBooking struct {
	Base
	UserID uuid.UUID 
	ClientID uuid.UUID 
	ConferenceID uuid.UUID 
	Source string
	Amount float64
	Discount float64
	DiscountType string
	AmountDue float64
	PaymentStatus string
	Currency string
	
}