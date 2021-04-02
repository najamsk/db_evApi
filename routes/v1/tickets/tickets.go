package favorites

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
	"strconv"

	//"github.com/satori/go.uuid"
)

var enviroment string = "staging"


var apiConfig *config.Config
var apiRouter *gin.Engine
var userService services.User
var ticketTypeService services.TicketType
var imageService services.Image
var paymentService services.Payment
var conferenceService services.Conference

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration
	apiRouter = router
	if(apiConfig.Items.Environment == "production"){
		enviroment = "live";
	}
	api := router.Group("api/v1/tickets")
	{
		api.POST("/detail", utils.BasicAuth(configuration), getTicketAndBookingInfo)
		api.POST("/type", utils.BasicAuth(configuration), getTicketInfo)
		api.POST("/book", utils.BasicAuth(configuration), bookTicket)
		api.POST("/user/bookings", utils.BasicAuth(configuration), userBookedTickets)
		api.POST("/consume", utils.BasicAuth(configuration), consumeTicket)	
	}
	apiv2 := router.Group("api/v2/tickets")
	{
		
		apiv2.POST("/detail", utils.BasicAuthV2(configuration), getTicketAndBookingInfo)
		apiv2.POST("/type", utils.BasicAuthV2(configuration), getTicketInfo)
		apiv2.POST("/book", utils.BasicAuthV2(configuration), bookTicket)
		apiv2.POST("/user/bookings", utils.BasicAuthV2(configuration), userBookedTickets)	
		apiv2.POST("/consume", utils.BasicAuthV2(configuration), consumeTicket)
		apiv2.POST("/reservation/cancel", utils.BasicAuthV2(configuration), cancelTicketReservation)
	}
}

func getTicketInfo(c *gin.Context) {
	var errorModel *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	var req viewmodels.Request;
	var ticketType []models.TicketType
	var bookings []viewmodels.TicketBookedVM
	var empty_message string
	type PaymentOptions struct {ID uuid.UUID; GatewayName string; GatewayType string}
	
	type Profile struct {ID uuid.UUID; FirstName string; LastName string; Email string; Organization string; Designation string; ProfileImage string}
	var user_profile Profile;
	c.BindJSON(&req)
	fmt.Println("error:", errorModel)
		var userdb, err = userService.Get(req.UserID);
		errorModel = utils.ErrorHandler(errorModel, err, "User does not exist.", "", "getTicketInfo")

		if(userdb == nil){
			errorModel = utils.AddError(errorModel, "User does not exist.", "User does not exist.", "getTicketInfo")
		}else{
			//check if user allow to scan QR code
			var isinRole = utils.IsInRole(userdb.Roles, apiConfig.Items.TicketSeller.Roles);
			if(isinRole){
			var memberdb, err = userService.Get(req.MemberID);
			errorModel = utils.ErrorHandler(errorModel, err, "Invalid QR code.", "", "getTicketInfo")
			if(memberdb == nil){
				errorModel = utils.AddError(errorModel, "Invalid QR code.", "User with provided MemberID not found.", "getTicketInfo")
			}else{
				user_profile.ID = memberdb.ID
				user_profile.FirstName = memberdb.FirstName
				user_profile.LastName = memberdb.LastName
				user_profile.Email = memberdb.Email
				user_profile.Organization = memberdb.Organization
				user_profile.Designation = memberdb.Designation
				

				var img, err = imageService.GetImage(memberdb.ID, "user", "user_profile");
				if(img != nil && err == nil){
					user_profile.ProfileImage = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
				}

				ticketType, err = ticketTypeService.GetAvailableTicketTypes(req.ConferenceID, req.ClientID);
				errorModel = utils.ErrorHandler(errorModel, err, "Error while getting ticket types.", "", "getTicketInfo")
				fmt.Println("ticketType:", ticketType)
				if(ticketType==nil || len(ticketType)<1){
					empty_message = "No ticket available."
					errorModel = utils.AddError(errorModel, "Ticket type(s) not found.", "Ticket type(s) not found.", "getTicketInfo")
				}
				bookings, err =  ticketTypeService.GetTicketTypeBooking(req.MemberID, req.ConferenceID);
			}
			}else{
				errorModel = utils.AddError(errorModel, "You are not authorized to get ticket info.", "You are not authorized to get ticket info.", "getTicketInfo")
			}
		}
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":errorModel, 
			"success_message": success_message, 
			"user_profile": user_profile,
			"ticketType":ticketType,
			"bookings": bookings,
			"empty_message":empty_message,
		},
	)
}

func bookTicket(c *gin.Context) {
	var errorModel *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var ticketId *uuid.UUID
	loggy.OpenLog()
	var conferencedb models.Conference;
	
	var ticketDB *models.Ticket
	//var bookings []viewmodels.TicketBookedVM
	
	type TicketInfo struct {UserID uuid.UUID; MemberID uuid.UUID; AmountPaid float64; TicketTypeID uuid.UUID; ConferenceID uuid.UUID; ClientID uuid.UUID}
	var req TicketInfo;
	c.BindJSON(&req)
	fmt.Println("error:", req)
		var userdb, err = userService.Get(req.UserID);
		errorModel = utils.ErrorHandler(errorModel, err, "User does not exist.", "", "bookTicket")

		if(userdb == nil){
			errorModel = utils.AddError(errorModel, "User does not exist.", "User does not exist.", "bookTicket")
		}else{
			//check if user allow to scan QR code
			var isinRole = utils.IsInRole(userdb.Roles, apiConfig.Items.TicketSeller.Roles);
			if(isinRole){
			var memberdb, err = userService.Get(req.MemberID);
			errorModel = utils.ErrorHandler(errorModel, err, "Attendee not found.", "", "bookTicket")
			if(memberdb == nil){
				errorModel = utils.AddError(errorModel, "Attendee not found.", "User with provided MemberID not found.", "bookTicket")
			}else{

				ticketId, err = ticketTypeService.BookTicket(req.UserID, req.ConferenceID, req.ClientID, req.MemberID, req.AmountPaid, req.TicketTypeID)
				errorModel = utils.ErrorHandler(errorModel, err, "Ticket not updated please try again.", "", "bookTicket")
				fmt.Println("ticketId:", ticketId)
				if(ticketId != nil){
				ticketDB, err = ticketTypeService.GetTicket(ticketId);
				fmt.Println("err:", err)
				errorModel = utils.ErrorHandler(errorModel, err, "Error while getting ticker after update. Please contact support team.", "", "bookTicket")
				if(ticketDB != nil){
					success_message = "Ticket "+ticketDB.SerialNo + " is successfully booked for "+memberdb.FirstName + " " + memberdb.LastName;
					//send email

				contactdb, _ := conferenceService.GetContact(req.ConferenceID);
				conferencedb, _ = conferenceService.Get(req.ConferenceID);
				utils.SendTicketPaymentConfirmationEmail(memberdb, conferencedb, req.ClientID, req.AmountPaid, "PKR", contactdb,  ticketDB.SerialNo);

				}
				}else{
					errorModel = utils.AddError(errorModel, "Ticket not found.", "Ticket not found.", "bookTicket")
				}
			}
			}else{
				errorModel = utils.AddError(errorModel, "You are not authorized to get ticket info.", "You are not authorized to get ticket info.", "bookTicket")
			}
		}
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":errorModel, 
			"success_message": success_message, 
			"ticketId":ticketId,
			"ticket":ticketDB,
		},
	)
}

func userBookedTickets(c *gin.Context) {
	var errorModel *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	var req viewmodels.Request;
	var bookings []viewmodels.TicketBookedVM
	var empty_message string
	
	type Profile struct {ID uuid.UUID; FirstName string; LastName string; Email string; Organization string; Designation string; ProfileImage string}
	var user_profile Profile;
	c.BindJSON(&req)
	fmt.Println("error:", errorModel)
		var userdb, err = userService.Get(req.UserID);
		errorModel = utils.ErrorHandler(errorModel, err, "User does not exist.", "", "userBookedTickets")

		if(userdb == nil){
			errorModel = utils.AddError(errorModel, "User does not exist.", "User does not exist.", "userBookedTickets")
		}else{
			//check if user allow to scan QR code
			var isinRole = utils.IsInRole(userdb.Roles, apiConfig.Items.TicketChecker.Roles);
			if(isinRole){
			var memberdb, err = userService.Get(req.MemberID);
			errorModel = utils.ErrorHandler(errorModel, err, "Invalid QR code.", "", "userBookedTickets")
			if(memberdb == nil){
				errorModel = utils.AddError(errorModel, "Invalid QR code.", "User with provided MemberID not found.", "userBookedTickets")
			}else{
				user_profile.ID = memberdb.ID
				user_profile.FirstName = memberdb.FirstName
				user_profile.LastName = memberdb.LastName
				user_profile.Email = memberdb.Email
				user_profile.Organization = memberdb.Organization
				user_profile.Designation = memberdb.Designation
				

				var img, err = imageService.GetImage(memberdb.ID, "user", "user_profile");
				if(img != nil && err == nil){
					user_profile.ProfileImage = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
				}

				bookings, err =  ticketTypeService.GetTicketTypeBooking(req.MemberID, req.ConferenceID);
				if(len(bookings)<1){
					var ticketStats = ticketTypeService.GetUserTicketStat(req.MemberID, req.ConferenceID);
					       
					if(ticketStats.ConsumedCount>0){
						empty_message = memberdb.FirstName + " "+memberdb.LastName +" ticket detail: \n";
						empty_message = empty_message + strconv.FormatInt(int64(ticketStats.ConsumedCount), 10) + " ticket(s) consumed."
					}
					 if(ticketStats.InactiveCount>0){
						empty_message = empty_message +"\n"+ strconv.FormatInt(int64(ticketStats.InactiveCount), 10) + " ticket(s) inactive."
					}
					 if(ticketStats.ExpiredCount>0){
						empty_message = empty_message +"\n"+ strconv.FormatInt(int64(ticketStats.ExpiredCount), 10)  + " ticket(s) expired."
					}

					empty_message = empty_message +"\n No ticket available to consume.";
				}
			}
			}else{
				errorModel = utils.AddError(errorModel, "You are not authorized to get ticket info.", "You are not authorized to get ticket info.", "userBookedTickets")
			}
		}
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":errorModel, 
			"success_message": success_message, 
			"user_profile": user_profile,
			"bookings": bookings,
			"empty_message":empty_message,
		},
	)
}

func consumeTicket(c *gin.Context) {
	var errorModel *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var consumedCount int
	loggy.OpenLog()
	
	
	type TicketInfo struct {UserID uuid.UUID; MemberID uuid.UUID; ConferenceID uuid.UUID; ClientID uuid.UUID}
	var req TicketInfo;
	c.BindJSON(&req)
	fmt.Println("error:", req)
		var userdb, err = userService.Get(req.UserID);
		errorModel = utils.ErrorHandler(errorModel, err, "User does not exist.", "", "consumeTicket")

		if(userdb == nil){
			errorModel = utils.AddError(errorModel, "User does not exist.", "User does not exist.", "consumeTicket")
		}else{
			//check if user allow to scan QR code
			var isinRole = utils.IsInRole(userdb.Roles, apiConfig.Items.TicketChecker.Roles);
			if(isinRole){
			var memberdb, err = userService.Get(req.MemberID);
			errorModel = utils.ErrorHandler(errorModel, err, "Attendee not found.", "", "consumeTicket")
			if(memberdb == nil){
				errorModel = utils.AddError(errorModel, "Attendee not found.", "User with provided MemberID not found.", "consumeTicket")
			}else{

				consumedCount, err = ticketTypeService.ConsumeTicket(req.UserID, req.ConferenceID, req.ClientID, req.MemberID)
				if(consumedCount>0){
					success_message = strconv.FormatInt(int64(consumedCount), 10)  + " ticket(s) consumed successfully."
				}else{ success_message = "Looks you do not have ticket(s) to consume."}
				errorModel = utils.ErrorHandler(errorModel, err, "Ticket not updated please try again.", "", "consumeTicket")
				fmt.Println("consumedCount:", consumedCount)
			}
			}else{
				errorModel = utils.AddError(errorModel, "You are not authorized to get ticket info.", "You are not authorized to get ticket info.", "consumeTicket")
			}
		}
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":errorModel, 
			"success_message": success_message, 
			"consumedCount":consumedCount,
		},
	)
}

func cancelTicketReservation(c *gin.Context) {
	var empty_message string
	var errorModel *viewmodels.Error
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	defer loggy.CloseLog()
	//type Request struct {PaymentIntegrationId uuid.UUID; TicketTypeId uuid.UUID; TicketCount int; Amount int64; CardNumber string; ExpiryMonth string; ExpiryYear string; CVC string;}
	var req viewmodels.TicketPaymentRequest;
	var success_message = "";
			
	c.BindJSON(&req)
	fmt.Println("req:", req);

	var result = ticketTypeService.CancelReserveTickets(uuid.FromStringOrNil(req.UserID), uuid.FromStringOrNil(req.ClientID), uuid.FromStringOrNil(req.ConferenceID), []string{}, apiConfig.Items.TicketReservationReleaseTimeInSec);

	c.JSON(
		status_code,
		gin.H{
			"error":errorModel, 
			"empty_message":empty_message,
			"success_message": success_message,
			"reservation_cancel_count": result,
		},
	)
}

func getTicketAndBookingInfo(c *gin.Context) {
	var errorModel *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	var req viewmodels.Request;
	var ticketType []models.TicketType
	//var bookings []viewmodels.TicketBookedVM
	var userTickets []viewmodels.TicketVM
	var empty_message string
	type PaymentOptions struct {ID uuid.UUID; GatewayName string; GatewayType string}
	var paymentOptions []PaymentOptions
	var paymentOptionsDB []models.PaymentIntegration
	var conferencedb models.Conference;
	var conference_title string;
	var conference_start_date time.Time;
	var conference_end_date time.Time;
	
	c.BindJSON(&req)
	fmt.Println("error:", errorModel)
		var userdb, err = userService.Get(req.UserID);
		errorModel = utils.ErrorHandler(errorModel, err, "User does not exist.", "", "getTicketAndBookingInfo")

		if(userdb == nil){
			errorModel = utils.AddError(errorModel, "User does not exist.", "User does not exist.", "getTicketAndBookingInfo")
		}else{
			conferencedb, err = conferenceService.Get(req.ConferenceID);
			if(err!=nil){
			 errorModel = utils.ErrorHandler(errorModel, err, "Problem with conference..", "", "makeTicketPayment")
			}else{
				conference_title = conferencedb.Title;
				conference_start_date = conferencedb.StartDate;
				conference_end_date = conferencedb.EndDate;
			}
			ticketType, err = ticketTypeService.GetAvailableTicketTypes(req.ConferenceID, req.ClientID);
			errorModel = utils.ErrorHandler(errorModel, err, "Error while getting ticket types.", "", "getTicketAndBookingInfo")
			fmt.Println("ticketType:", ticketType)
			if(ticketType==nil || len(ticketType)<1){
					empty_message = "No ticket available."
					//errorModel = utils.AddError(errorModel, "Ticket type(s) not found.", "Ticket type(s) not found.", "getTicketAndBookingInfo")
				}
		//	bookings, err =  ticketTypeService.GetTicketTypeBooking(req.UserID, req.ConferenceID);
			userTickets, err = ticketTypeService.GetUserAllTicketsByConference(req.UserID, req.ConferenceID);
			errorModel = utils.ErrorHandler(errorModel, err, "Error while getting ticket bookings.", "", "getTicketAndBookingInfo")

			

	
	paymentOptionsDB, err = paymentService.GetPaymentGatewayOptions(req.ClientID, enviroment);
	errorModel = utils.ErrorHandler(errorModel, err, "Payment Gateway down.", "", "getTicketAndBookingInfo")

 if(len(paymentOptionsDB)>0){
	for _, gateway := range paymentOptionsDB {
		paymentOptions = append(paymentOptions, PaymentOptions{ID:gateway.ID, GatewayName:gateway.GatewayName, GatewayType: gateway.GatewayType})
	}

	}

		}
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":errorModel, 
			"success_message": success_message, 
			"ticketType":ticketType,
			"user_tickets": userTickets,
			"empty_message":empty_message,
			"paymentOptions":paymentOptions,
			"conference_title": conference_title,
		    "conference_start_date": conference_start_date,
			"conference_end_date": conference_end_date,
		},
	)
}







