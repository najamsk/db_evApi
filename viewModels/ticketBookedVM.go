package viewmodels

import(
	"time"
)

type TicketBookedVM struct {
	TicketType    	string
	BookedCount     string
	UnitPrice      float64
	TotalPrice     float64
	Currency		string
	ValidFrom			time.Time
	ValidTo			time.Time			
}

