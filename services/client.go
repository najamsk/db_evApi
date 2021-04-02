package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"fmt"
	"github.com/satori/go.uuid"
)

type Client struct {}


func (client Client) Get(id uuid.UUID) (models.Client, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var clientdb models.Client
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	var err = db.Find(&clientdb, "ID = ?", id).Error
	fmt.Println(err)
	if err != nil{
		return clientdb, err
	}
	
	return clientdb, nil
}