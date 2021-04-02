package payments

import (
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/gin-gonic/gin"
	"github.com/najamsk/eventvisor/eventvisor.api/services"
	"github.com/satori/go.uuid"
	"fmt"
	"time"
	// "strconv"
	 "strings"
	"errors"
)

type PaymentResult struct {
	Status string
	AmountPaid float64
	TransactionID string
	SuccessMessage string
	AlertMessage string
	PaymentGateway string
	BookedTickets []string
   }

var apiConfig *config.Config
var apiRouter *gin.Engine
var paymentService services.Payment
var ticketTypeService services.TicketType
var ticketBookingService services.TicketBooking
var ticketBookingItemService services.TBookingItem
var conferenceService services.Conference
var userService services.User
var enviroment string = "staging"

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration
	apiRouter = router

	if(apiConfig.Items.Environment == "production"){
		enviroment = "live";
	}
	
	apiv2 := router.Group("api/v2/payments")
	{
		apiv2.POST("/gateways", utils.BasicAuthV2(configuration), getPaymentGateWaysOptions)
		apiv2.POST("/ticket/paynow", utils.BasicAuthV2(configuration), makeTicketPayment)
	}
}

func getPaymentGateWaysOptions(c *gin.Context) {
	var err = error(nil);
	var empty_message string
	var errorModel *viewmodels.Error
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	defer loggy.CloseLog()
	var req viewmodels.Request;
	type PaymentOptions struct {ID uuid.UUID; GatewayName string; GatewayType string}
	
	c.BindJSON(&req)
	
	var paymentOptions []PaymentOptions
	var paymentOptionsDB []models.PaymentIntegration

	paymentOptionsDB, err = paymentService.GetPaymentGatewayOptions(req.ClientID, enviroment);
	errorModel = utils.ErrorHandler(errorModel, err, "Payment Gateway down.", "", "getPaymentGateWaysOptions")
	if(len(paymentOptionsDB)>0){
		for _, gateway := range paymentOptionsDB {
			paymentOptions = append(paymentOptions, PaymentOptions{ID:gateway.ID, GatewayName:gateway.GatewayName, GatewayType: gateway.GatewayType})
		}
	
		}

	c.JSON(
		status_code,
		gin.H{
			"error":errorModel, 
			"empty_message":empty_message,
			"paymentOptions": paymentOptions,
		},
	)
}

func makeTicketPayment(c *gin.Context) {
	var err = error(nil);
	var empty_message string
	var errorModel *viewmodels.Error
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	defer loggy.CloseLog()
	var req viewmodels.TicketPaymentRequest;
	var success_message = "";
	var ticketIds []string
	var conferencedb models.Conference;
	var conference_title string;
	var conference_start_date time.Time;
	var conference_end_date time.Time;
	var paymentResult PaymentResult;
	var userdb *models.User
	var currency string
	
	err = c.BindJSON(&req)
	fmt.Println("req.err:", err);
	fmt.Println("req:", req);

	errorModel = utils.ValidateTicketPaymentRequest(req);
	conferencedb, err = conferenceService.Get(uuid.FromStringOrNil(req.ConferenceID));
	if(err!=nil){
			errorModel = utils.ErrorHandler(errorModel, err, "Problem with conference.", "", "makeTicketPayment")
		}
	userdb, err = userService.Get(uuid.FromStringOrNil(req.UserID));
	if(err!=nil){
		errorModel = utils.ErrorHandler(errorModel, err, "Problem with user account.", "", "makeTicketPayment")
	}
	if(errorModel == nil){

		conference_title = conferencedb.Title;
		conference_start_date = conferencedb.StartDate;
		conference_end_date = conferencedb.EndDate;

		ticketIds, errorModel = reserveTicket(req, errorModel);
		//
		if(errorModel == nil){

	var paymentOptionsDB = new(models.PaymentIntegration)
	paymentOptionsDB, err = paymentService.Get(uuid.FromStringOrNil(req.PaymentIntegrationID));

	if(err!=nil){
		cacceledTicket := ticketTypeService.CancelReserveTickets(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), ticketIds, 0);
		errorModel = utils.ErrorHandler(errorModel, err, "Sorry, payment service is down.", "", "makeTicketPayment")
		fmt.Println("cacceledTicket:", cacceledTicket);
	}else if(paymentOptionsDB == nil){
		cacceledTicket := ticketTypeService.CancelReserveTickets(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), ticketIds, 0);
		errorModel = utils.AddError(errorModel, "Sorry, payment service is down.", "", "makeTicketPayment")
		fmt.Println("cacceledTicket:", cacceledTicket);
		}	
	if(errorModel == nil && paymentOptionsDB != nil){
		var ticketBooking models.TicketBooking
		ticketBooking, errorModel = InsertBookingDB(req, errorModel);

		if(errorModel == nil && ticketBooking.AmountDue > 0){
			if(paymentOptionsDB.GatewayType == "stripe"){
				currency = ticketBooking.Currency;
				paymentResult, errorModel = MakeStripePayment(req,  paymentOptionsDB, ticketBooking,  ticketIds, errorModel);
				//isPaymentSuccessful, amountPaid, errorModel = MakeStripePayment(req,  paymentOptionsDB, ticketBooking,  ticketIds, errorModel);
				
				if(paymentResult.Status != "succeeded"){
					cacceledTicket := ticketTypeService.CancelReserveTickets(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), ticketIds, 0);
					fmt.Println("cacceledTicket:", cacceledTicket);
				}

				if(paymentResult.Status == "succeeded" && len(paymentResult.BookedTickets)>0){
					contactdb, _ := conferenceService.GetContact(conferencedb.ID);
					var tickets string = strings.Join(paymentResult.BookedTickets,", ");
					utils.SendTicketPaymentConfirmationEmail(userdb, conferencedb, uuid.FromStringOrNil(req.ClientID), paymentResult.AmountPaid, ticketBooking.Currency, contactdb,  tickets);
				}
		   }
		}else{
			cacceledTicket := ticketTypeService.CancelReserveTickets(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), ticketIds, 0);
			fmt.Println("cacceledTicket:", cacceledTicket);
		}	
	}
}
	}
	c.JSON(
		status_code,
		gin.H{
			"error":errorModel, 
			"empty_message":empty_message,
			"success_message": success_message,
			"amount_paid": paymentResult.AmountPaid,
			"conference_title": conference_title,
		    "conference_start_date": conference_start_date,
			"conference_end_date": conference_end_date,
			"alert_message": paymentResult.AlertMessage,
			"currency": currency,
		},
	)
}

//check if ticket available from database
//if available reserve tickets and return ticketids reserved
//if reserve tickets count not eqal to ticket required then cancel/release reserved ticket 
func reserveTicket(req viewmodels.TicketPaymentRequest, errorModel *viewmodels.Error) ([]string, *viewmodels.Error){

	var ticket_available = ticketTypeService.CheckTicketAvailability(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), req.Tickets);
	var reserve_ticketIds []string
	var ticker_required int
	if(ticket_available){
		reserve_ticketIds, ticker_required = ticketTypeService.ReserveTickets(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), req.Tickets, apiConfig.Items.TicketExpiryInSec);
		if(len(reserve_ticketIds) != ticker_required && ticker_required >0){
			cacceledTicket := ticketTypeService.CancelReserveTickets(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), reserve_ticketIds, 0);
			fmt.Println("cacceledTicket:", cacceledTicket);
			if(errorModel == nil){
				errorModel = new(viewmodels.Error)
				}
				fmt.Println("errorModel:",errorModel);
				errorModel.DisplayErrors = append(errorModel.DisplayErrors, "Ticket(s) not available.")
				errorModel.InnerErrors = append(errorModel.InnerErrors, "Ticket(s) not available.")
				errorModel.ApiStatusCode = http.StatusBadRequest
		}
	}else{
		if(errorModel == nil){
			errorModel = new(viewmodels.Error)
			}
			fmt.Println("errorModel:",errorModel);
			errorModel.DisplayErrors = append(errorModel.DisplayErrors, "Ticket(s) not available.")
			errorModel.InnerErrors = append(errorModel.InnerErrors, "Ticket(s) not available.")
			errorModel.ApiStatusCode = http.StatusBadRequest
	}
		
  return reserve_ticketIds, errorModel;
}

//get tickettypeids and quantity from request and get tickettype data from db
//insert ticketbooking and ticketboolingitems in database
//return type ticketbooling model and error model

func InsertBookingDB(req viewmodels.TicketPaymentRequest, errorModel *viewmodels.Error) (models.TicketBooking, *viewmodels.Error){
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	defer loggy.CloseLog()
	var err = error(nil);
	var ticketBooking models.TicketBooking
	var toalAmount float64 =0;
	var currency string;
	var ticketIds []string

	//populate ticketids array and map tickettypeid as key and quantity as value 
	//so that we can get quantity against tickettypeid while calculating price
		tTypesMap := make(map[string]int)
		for _, ticketType := range req.Tickets {
			ticketIds = append(ticketIds, ticketType.TicketTypeId)
			tTypesMap[ticketType.TicketTypeId] = ticketType.Quantity
		}
		var ticketsType []*models.TicketType
		ticketsType, err = ticketTypeService.GetTicketTypes(ticketIds);

		if(err!=nil){
			errorModel = utils.ErrorHandler(errorModel, err, "Error while booking ticket(s).", "", "InsertBookingDB")
		}
		
		fmt.Println("ticketsType:", ticketsType);
		for _, ticketTypedb := range ticketsType {
			toalAmount =toalAmount + (ticketTypedb.Amount*float64(tTypesMap[ticketTypedb.ID.String()]))
			currency = ticketTypedb.AmmountCurrency;
		}
		
		if(toalAmount>0){
		//save booking in db
		ticketBooking.UserID = uuid.FromStringOrNil(req.UserID)
		ticketBooking.ClientID = uuid.FromStringOrNil(req.ClientID)
		ticketBooking.ConferenceID = uuid.FromStringOrNil(req.ConferenceID)
		ticketBooking.Source = "api";
		ticketBooking.Amount = toalAmount;
		ticketBooking.AmountDue = toalAmount;
		ticketBooking.PaymentStatus = "pending"
		ticketBooking.Currency = currency
		ticketBooking, err = ticketBookingService.Insert(ticketBooking)

		if(err!=nil){
			errorModel = utils.ErrorHandler(errorModel, err, "Error while booking ticket(s).", "", "InsertBookingDB")
		}else{
				// save booking items in db
				for _, ticketTypedb := range ticketsType {
					
					var ticketBookingItem models.TicketBookingItem
					ticketBookingItem.TicketBookingID = ticketBooking.ID
					ticketBookingItem.TicketTypeID = ticketTypedb.ID
					ticketBookingItem.Quantity = tTypesMap[ticketTypedb.ID.String()];
					ticketBookingItem.UnitPrice = ticketTypedb.Amount
					ticketBookingItem.TotalPrice = ticketTypedb.Amount*float64(tTypesMap[ticketTypedb.ID.String()])
					ticketBookingItem.AmountDue = ticketBookingItem.TotalPrice;
					ticketBookingItem.Currency = ticketTypedb.AmmountCurrency;
					ticketBookingItem.PaymentStatus = "pending";
					ticketBookingItem, err = ticketBookingItemService.Insert(ticketBookingItem);
					if(err!=nil){
						errorModel = utils.ErrorHandler(errorModel, err, "Error while booking ticket(s).", "", "InsertBookingDB")
					}
				}
		}
		}else{
				err = errors.New("Error while booking ticket(s).");
				errorModel = utils.ErrorHandler(errorModel, err, "Error while booking ticket(s).", "", "InsertBookingDB");
		} 
		return ticketBooking, errorModel
}

//update booking status from pending to paid in db
//update boolingitems status from pending to paid in db
//update reserved tickets booked by field
func updateBooking(req viewmodels.TicketPaymentRequest, bookingId uuid.UUID, ticketIds []string, errorModel *viewmodels.Error) ([]string, bool, *viewmodels.Error){

	//var updated_records int =0 
	var err = error(nil);
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	defer loggy.CloseLog()

	var isUpdated bool = true;
	var booked_tickets []string;

	err = ticketBookingService.UpdateBookingStatus(bookingId);
    fmt.Println("err:",err);
	
	if(err != nil){
		isUpdated = false;
		loggy.Logger.Info().Msg("updateBooking.UpdateBookingStatus.Error: "+err.Error())
		//errorModel = utils.ErrorHandler(errorModel, err, "Error while updating ticket booking status.", "", "updateBooking")
	} 

	err = ticketBookingService.UpdateBookingItemStatus(bookingId);
	fmt.Println("err:",err);
	if(err != nil){
		isUpdated = false;
		loggy.Logger.Info().Msg("updateBooking.UpdateBookingItemStatus.Error: "+err.Error())
		//errorModel = utils.ErrorHandler(errorModel, err, "Error while updating ticket booking status.", "", "updateBooking")
	} 
	
	 booked_tickets, err = ticketTypeService.UpdateTicketBooking(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), ticketIds);

	if(err != nil){
		isUpdated = false;
		loggy.Logger.Info().Msg("updateBooking.UpdateTicketBooking.Error: "+err.Error())
			//errorModel = utils.ErrorHandler(errorModel, err, "Error while updating ticket booking status.", "", "updateBooking")
		} 
	

		return booked_tickets, isUpdated, errorModel
}

//Make payment on stripe 
//if payment successful on stripe write response in log and in database 
//update booking status , booling item status and ticket booked by column
//insert payment in database

func MakeStripePayment(req viewmodels.TicketPaymentRequest, paymentOptionsDB *models.PaymentIntegration, ticketBooking models.TicketBooking, ticketIds []string, errorModel *viewmodels.Error) (PaymentResult, *viewmodels.Error){
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	defer loggy.CloseLog()

	var err = error(nil);
	var paymentResult PaymentResult

	res_stripe, errorModel := utils.MakePaymentCreditCardStripe(paymentOptionsDB.SecretKey, req.CardNumber, req.ExpiryMonth, req.ExpiryYear, req.CVC, ticketBooking.AmountDue, ticketBooking.Currency, ticketBooking.ID.String());
	//fmt.Println("payment_response_stripe:",res_stripe)
	if(errorModel == nil && res_stripe != nil){

		paymentResult.AmountPaid = float64(res_stripe.Amount)/100;
		paymentResult.TransactionID = res_stripe.ID;
		paymentResult.Status = res_stripe.Status;
		paymentResult.PaymentGateway = "stripe";

		loggy.Logger.Info().Msg("MakeStripePayment:StripeResponse: "+utils.ConvertToJSON(res_stripe));
		var paymentRes models.PaymentResponse;
		paymentRes.PaymentIntegrationID = uuid.FromStringOrNil(req.PaymentIntegrationID);
		paymentRes.GatewayName = paymentOptionsDB.GatewayType;
		paymentRes.BookingID = ticketBooking.ID;
		paymentRes.Request = utils.ConvertToJSON(req);
		paymentRes.Response = utils.ConvertToJSON(res_stripe);
		paymentRes.Environment = enviroment;
		paymentRes, err = paymentService.InsertPaymentResponse(paymentRes);
		if(err!=nil){
			loggy.Logger.Info().Msg("MakeStripePayment:InsertPaymentResponse: "+err.Error());
		}
		
		//update booking status , booling item status and ticket booked by column
		var isBookingUpdated bool = false;
		paymentResult.BookedTickets, isBookingUpdated, errorModel = updateBooking(req, ticketBooking.ID, ticketIds, errorModel);
		
		var paymentMethod string;
		var sourceType string;
		var status string;
		var sourceBrand string;
		var sourcelast4 string;

		if(res_stripe.Status == "succeeded"){
			status = "paid";
		}
		if(res_stripe.PaymentMethodDetails != nil){

			paymentMethod = string(res_stripe.PaymentMethodDetails.Type);
			sourceType = string(res_stripe.PaymentMethodDetails.Type);

			if(res_stripe.PaymentMethodDetails.Card != nil){
				sourcelast4 = res_stripe.PaymentMethodDetails.Card.Last4;
				sourceBrand = string(res_stripe.PaymentMethodDetails.Card.Brand);
			}
		}

		var isPaymentUpdated bool = false;
		isPaymentUpdated, errorModel = InsertDBPayment(enviroment, ticketBooking.ID, "ticket_booking", req, 
									ticketBooking.AmountDue, paymentResult.AmountPaid, string(res_stripe.Currency), res_stripe.ID, paymentMethod,
									sourceType, sourcelast4, sourceBrand, status, errorModel);

			if(!(isPaymentUpdated && isBookingUpdated)){
				amount  := fmt.Sprintf("%f", paymentResult.AmountPaid)
				amount = strings.TrimRight(strings.TrimRight(amount, "0"), ".")
				paymentResult.AlertMessage = amount + " ("+string(res_stripe.Currency)+") has been deducted from your account but there is an issue while updating ticket booking in our system. Please contact customer support team.";
				loggy.Logger.Info().Msg("MakeStripePayment:Error: "+paymentResult.AlertMessage);
				//email while issue updating data
				//emailBody := "Stripe payment is done but issue while updating booking, tickets or payment. "
				//emailBody = emailBody+"BookingID: "+ticketBooking.ID.String();
				//tils.SendEmail(apiConfig.Items.MoftakDevEmails,[]string{},"Issue in ticket booking",emailBody, "",apiConfig);
				}
			}
	return paymentResult, errorModel;
}

func InsertDBPayment(enviroment string, entityID uuid.UUID, entityType string, req viewmodels.TicketPaymentRequest, amountDue float64, amountPaid float64, currency string, transactionID string, paymentMethod string, sourceType string, sourceLast4 string, sourceBrand string, status string,   errorModel *viewmodels.Error) (bool, *viewmodels.Error){
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	defer loggy.CloseLog()

	var isUpdated bool = true;

	var err = error(nil);
	var payment models.Payment
		payment.Environment = enviroment;
		payment.EntityID = entityID; 
		payment.EntityType = entityType; 
		payment.PaymentFromUserID = uuid.FromStringOrNil(req.UserID);
		payment.PaymentIntegrationID = uuid.FromStringOrNil(req.PaymentIntegrationID);
		payment.ClientID = uuid.FromStringOrNil(req.ClientID);
		payment.ConferenceID = uuid.FromStringOrNil(req.ConferenceID);
		payment.AmountDue = amountDue;
		payment.AmountPaid = amountPaid;
		payment.Currency = currency;
		payment.TransactionID = transactionID;
		payment.Status = status;
		payment.PaymentMethod = paymentMethod; 
		payment.SourceType = sourceType; 
		payment.SourceLast4 = sourceLast4;
		payment.SourceBrand = sourceBrand;
		
		payment, err = paymentService.Insert(payment);

		if(err != nil){
			isUpdated = false;
			loggy.Logger.Info().Msg("InsertDBPayment.paymentService.Insert.Error: "+err.Error())
			//errorModel = utils.ErrorHandler(errorModel, err, "Error while updating payment.", "", "InsertDBPayment")
		} 

		return isUpdated, errorModel;
}









