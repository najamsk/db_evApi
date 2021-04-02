package models

import(
	"github.com/satori/go.uuid"
)

type EmailTemplate struct {
	Base
	Name         	string
	Title         	string
	Subject			string
	EmailBody		string
	TemplateTypeID  uuid.UUID
	IsHtml     		bool
	IsDefault 		bool
	IsActive 		bool
	ClientID        uuid.UUID
	ConferenceID	uuid.UUID
	EmailFrom		string
}


 
