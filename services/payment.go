package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
    _"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	_"time"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Payment struct {}

func (payment *Payment) Insert(data models.Payment) (models.Payment, error) {
	db := utils.GetDb()

	err := db.Create(&data).Error
	
	if err != nil{
		return data, err
	}
	
	return data, nil
}

func (payment *Payment) InsertPaymentResponse(data models.PaymentResponse) (models.PaymentResponse, error) {
	db := utils.GetDb()

	err := db.Create(&data).Error
	
	if err != nil{
		return data, err
	}
	
	return data, nil
}

func (payment *Payment) Update(data *models.Payment) (error) {
	db := utils.GetDb()

	// If you only want to update changed Fields, you could use Update, Updates
	//var err = db.Update(&userobj).Error
	//fmt.Println("err update user:", err)
	// Save will include all fields when perform the Updating SQL, even it is not changed
	 var err = db.Save(&data).Error
	 
	return err
}

func (payment *Payment) Get(id uuid.UUID) (*models.PaymentIntegration, error) {
	db := utils.GetDb()
	
	var paymentGateway models.PaymentIntegration
	var err = db.Find(&paymentGateway, "ID = ?", id).Error
	
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return &paymentGateway, nil
}

func (payment *Payment) GetPaymentGatewayOptions(clientid uuid.UUID, enviroment string) ([]models.PaymentIntegration, error) {
	db := utils.GetDb()
	var paymentOptions []models.PaymentIntegration
	
	var err = db.Where("id in( select cpi.payment_integration_id from client_payment_integrations cpi where cpi.client_id = ? and cpi.is_active = ?) and payment_integrations.environment = ? ", clientid, true, enviroment).Find(&paymentOptions).Error;
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return paymentOptions, nil
}