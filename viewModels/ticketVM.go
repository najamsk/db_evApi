package viewmodels

import ("time" 
//"github.com/satori/go.uuid"
)

type TicketVM struct {
	//ID				uuid.UUID
	//IsActive        bool
	//ValidFrom       time.Time
	ValidTo         time.Time
	//ClientID        uuid.UUID
	//ConferenceID    uuid.UUID
	SerialNo        string 
	//BookedBy        uuid.UUID 
	IsConsumed      bool
	//ConsumedBy      uuid.UUID 
	//ConsumedAt      *time.Time
	//SoldBy          uuid.UUID
	//TicketTypeID    uuid.UUID
	//BookedAt      	*time.Time
	//ReservedBy		uuid.UUID
	//ReservedAt		*time.Time 
	Price      		float64
	IsExpire        bool
	Currency        string 
}