package models

import ("time" 
"github.com/satori/go.uuid")

type Subscription struct {
	Base
	IsActive         bool
	StartDate        time.Time
	StartDateDisplay string
	EndDate          time.Time
	EndDateDisplay   string
	DurationDisplay  string
	ClientID         uuid.UUID
	Billed           float64
	BilledCurrency   string
	Payments         []SubscriptionPayment
	Remarks          string

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}