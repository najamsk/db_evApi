package models

import(

	 "github.com/satori/go.uuid"
)

type TicketBookingItem struct {
	Base
	TicketBookingID uuid.UUID 
	TicketTypeID uuid.UUID 
	Quantity int
	UnitPrice float64
	TotalPrice float64
	Discount float64
	DiscountType string
	AmountDue float64
	PaymentStatus string
	Currency string
}