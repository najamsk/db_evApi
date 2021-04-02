package models

import(
	"github.com/satori/go.uuid"
)

type ClientSmtpSetting struct {
	Base
	Host         	string
	Port         	int
	UserName		string
	Password		string
	IsActive 		bool
	ClientID        uuid.UUID
	ConferenceID	uuid.UUID
	EmailFrom		string
}


 
