package models

import(
"github.com/jinzhu/gorm"
"github.com/satori/go.uuid"
"time"
_"fmt")

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"created_at"`
    UpdatedAt *time.Time  `json:"update_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedBy *uuid.UUID `gorm:"type:uuid;"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid;"`
	DeletedBy *uuid.UUID `gorm:"type:uuid;"`
	Deleted bool
   }

   /* BeforeCreate will set a UUID rather than numeric ID. 
		 If you want to user auto generated uuid as primary key 
		 then uncomment following code */
   

   func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
	 return err
	}
	return scope.SetColumn("ID", uuid)
   } 