package models

import "github.com/satori/go.uuid"

type SubscriptionPayment struct {
	Base
	IsActive        bool
	Ammount         float64
	AmmountCurrency string
	PaymentLog      string
	PaymentGateway  string
	SubscriptionID  uuid.UUID
	TransactionID string

	// CreatedAt       time.Time
	// UpdatedAt       time.Time
	// DeletedAt       time.Time
}