package models
import "github.com/satori/go.uuid"

type SessionTicket struct {
	Base
	TicketTypeID   uuid.UUID
	TicketIssued   int64 //starts with non zero number
	TicketConsumed int64 //starts with zero number and should reach to ticketissued
	SessionID      uuid.UUID
}