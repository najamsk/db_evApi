package models

import ("time" 
"github.com/satori/go.uuid")

type Ticket struct {
	Base
	IsActive        bool
	ValidFrom       time.Time
	ValidTo         time.Time
	ClientID        uuid.UUID
	ConferenceID    uuid.UUID
	SerialNo        string //uuid maybe?
	BookedBy        uuid.UUID    //who purchased this
	IsConsumed      bool
	ConsumedBy      uuid.UUID       //should be same person who booked
	ConsumedAt      *time.Time //defualt value should be nil until its consumed
	ConsumedSession int
	SoldBy          uuid.UUID
	TicketTypeID    uuid.UUID
	BookedAt      	*time.Time
	ReservedBy		uuid.UUID
	ReservedAt		*time.Time 
	ReservationExpireAt *time.Time
	AmountPaid      float64
}