package models

import(
	"github.com/satori/go.uuid"
	_"fmt"
"time"
"github.com/jinzhu/gorm")

type PaymentResponse struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	PaymentIntegrationID uuid.UUID
	GatewayName string 
	PaymentID *uuid.UUID 
	BookingID uuid.UUID
	Request string 
	Response string
	Environment string
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (paymentres *PaymentResponse) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
	 return err
	}
	return scope.SetColumn("ID", uuid)
   } 
