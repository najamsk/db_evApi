package viewmodels

import(
	//"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid")

type Request struct {
	ClientID		uuid.UUID
	UserID			uuid.UUID
	SpeakerID		uuid.UUID
	MACAddress		string
	Platform		string
	ConferenceID	uuid.UUID
	Favorite		bool
	Email			string
	Password		string
	ConfirmPassword string
	PasswordToken string
	EntityID		uuid.UUID
	EntityType		string
	StartupID 		uuid.UUID
	OldPassword		string
	MemberID			uuid.UUID
	PageNo		int
	OffSet		int
	SearchText string
}

type TicketPaymentRequest struct {
	
	ClientID		string
	UserID			string
	MACAddress		string
	Platform		string
	ConferenceID	string
	EntityID		string
	EntityType		string

	PaymentIntegrationID string
	CardNumber string
	ExpiryMonth string
	ExpiryYear string
	CVC string
	Tickets []Ticket
	//ID *uuid.UUID
}

type Ticket struct {
	TicketTypeId string 
	Quantity int
	Price int64
	TotalAmount int64
	Discount int
	DiscountType string
}