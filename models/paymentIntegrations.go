package models

import(
	"github.com/satori/go.uuid"
	_"fmt")

type PaymentIntegration struct {
	Base
	GatewayName string
	GatewayCategoryId uuid.UUID
	GatewayType string
	IsPlatform  bool
	Environment string
	IsActive bool
	Country string
	IsLocal bool
	Url	string
	SecretKey  string
	KeyPassword string
	UserName string
	Password string
	Currency string
}
