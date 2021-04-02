package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"fmt"
    _"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	_"time"
	"github.com/satori/go.uuid"
)
type Role struct {}
func (role Role) GetRoleByName(role_name string) (models.Role, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var roledb models.Role
	var err = db.Where("name = ?", role_name).First(&roledb).Error
	fmt.Println(err)
	if err != nil{
		return roledb, err
	}
	
	return roledb, nil
}

func (role Role) InsertUserRole(roleId uuid.UUID, userId uuid.UUID, clientId uuid.UUID) (error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	//var roledb models.Role
	var query string = "insert into users_roles(user_id, role_id, client_id) values(?, ?, ?)";
	var err = db.Exec(query, userId, roleId, clientId).Error;
	fmt.Println(err)
	return err
}

func (role Role) GetUserRoles(userId uuid.UUID, clientId uuid.UUID) ([]*models.Role, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var roles []*models.Role
	var query string = " select * from roles where id in(select role_id from users_roles where user_id = ? and client_id = ?)";
	var err = db.Raw(query, userId, clientId).Scan(&roles).Error
	
	fmt.Println(err)
	if err != nil{
		return roles, err
	}
	
	return roles, nil
}