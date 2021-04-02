package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"fmt"
	"github.com/satori/go.uuid"
)

type EmailTemplate struct {}


func (template EmailTemplate) Get(id uuid.UUID) (models.EmailTemplate, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var templateDB models.EmailTemplate
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	var err = db.Find(&templateDB, "ID = ?", id).Error
	fmt.Println(err)
	if err != nil{
		return templateDB, err
	}
	
	return templateDB, nil
}

func (template EmailTemplate) GetTemplate(clientId uuid.UUID, conferenceId uuid.UUID, templateName string) (models.EmailTemplate, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var templateDB models.EmailTemplate
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	var err = db.Find(&templateDB, "client_id = ? and conference_id = ? and name = ? and is_active = ?", clientId, conferenceId, templateName, true).Error
	fmt.Println(err)
	if err != nil{
		return templateDB, err
	}
	
	return templateDB, nil
}