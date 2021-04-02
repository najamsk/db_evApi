package models

import(
	"github.com/satori/go.uuid")

type Payment struct {
	Base
	ParentPaymentID	*uuid.UUID
	EntityID uuid.UUID
	EntityType string
	PaymentFromUserID uuid.UUID
	PaymentIntegrationID uuid.UUID
	ClientID uuid.UUID
	ConferenceID uuid.UUID
	AmountDue float64
	AmountPaid float64 
	AttemptCount int
	Currency string
	SourceBrand string
	SourceID string 
	SourceLast4 string
	SourceType string
	Status string
	Environment string
	TransactionID string
	PaymentMethod string
	PaymentType string
	Quantity int	
}