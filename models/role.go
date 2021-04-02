package models

// User type
// type Role struct {
// 	Base
// 	RoleName string
// 	LoweredRoleName string 
// 	Description string
// 	Users	[]User `gorm:"many2many:user_roles;"`
// }

type Role struct {
	Base
	Name      string
	DisplayName string
	Weight int
}

