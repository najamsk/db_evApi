package utils

import (
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"github.com/satori/go.uuid"
	//"net/http"
	"regexp"
	"strconv"
	"time"
	"fmt"
	"errors"
)

func ValidateTicketPaymentRequest(request viewmodels.TicketPaymentRequest) *viewmodels.Error {

	var err = error(nil);
	var numeric string = "^[-+]?[0-9]+$"
	var rxNumeric  = regexp.MustCompile(numeric)
	var errorModel *viewmodels.Error
	var expiryYear int;
	var currentYear int = time.Now().Year();
	var currentMonth int = int(time.Now().Month());
	
	//validate UserID
	if(uuid.FromStringOrNil(request.UserID) == uuid.Nil){
		err := errors.New("Please provide valid UserID.");
		if(err != nil){
			errorModel = ErrorHandler(errorModel, err, "Problem with user account.", "", "ValidateTicketPaymentRequest")
		} 
	}

	//validate ConferenceID
	if(uuid.FromStringOrNil(request.ConferenceID) == uuid.Nil){
		err := errors.New("Please provide valid ConferenceID.");
		if(err != nil){
			errorModel = ErrorHandler(errorModel, err, "Problem with conference.", "", "ValidateTicketPaymentRequest")
		}
	}

	//validate ClientID
	if(uuid.FromStringOrNil(request.ClientID) == uuid.Nil){
		err := errors.New("Please provide valid ClientID.");
		if(err != nil){
			errorModel = ErrorHandler(errorModel, err, "Bad Request.", "", "ValidateTicketPaymentRequest")
		}
	}

	//validate integrationid
	if(uuid.FromStringOrNil(request.PaymentIntegrationID) == uuid.Nil){
		err := errors.New("Please provide valid payment gateway.");
		if(err != nil){
			errorModel = ErrorHandler(errorModel, err, "Sorry, payment service is down.", "", "ValidateTicketPaymentRequest")
		}
	}

	//validate card number
	if(len(request.CardNumber) <1 || !rxNumeric.MatchString(request.CardNumber)){
		err := errors.New("Please provide valid card number.");
		if(err != nil){
			errorModel = ErrorHandler(errorModel, err, "Please provide valid card number.", "", "ValidateTicketPaymentRequest")
		}
	}

	//validate expity year
	if(len(request.ExpiryYear) <1 || !rxNumeric.MatchString(request.ExpiryYear)){
		err := errors.New("Please provide valid expiry year.");
		if(err != nil){
			errorModel = ErrorHandler(errorModel, err, "Please provide valid expiry year.", "", "ValidateTicketPaymentRequest")
		}
	}else{
		
		expiryYear, err = strconv.Atoi(request.ExpiryYear)
		if (err != nil || expiryYear < currentYear) {
			err := errors.New("Please provide valid expiry year.");
			if(err != nil){
				errorModel = ErrorHandler(errorModel, err, "Please provide valid expiry year.", "", "ValidateTicketPaymentRequest")
			}
		}
	}

	//validate expiry month
	if(len(request.ExpiryMonth) <1 || !rxNumeric.MatchString(request.ExpiryMonth)){
			err := errors.New("Please provide valid expiry month.");
			if(err != nil){
				errorModel = ErrorHandler(errorModel, err, "Please provide valid expiry month.", "", "ValidateTicketPaymentRequest")
			}
	}else{
			n, err := strconv.Atoi(request.ExpiryMonth)
			if (err != nil || n<1 || n>12 || (expiryYear == currentYear && n<currentMonth) ) {
				err := errors.New("Please provide valid expiry month.");
				if(err != nil){
					errorModel = ErrorHandler(errorModel, err, "Please provide valid expiry month.", "", "ValidateTicketPaymentRequest")
				}
			}
	}

	//validate cvc
	if(len(request.CVC) <3 || !rxNumeric.MatchString(request.CVC)){
		err := errors.New("Please provide valid CVC.");
		if(err != nil){
			errorModel = ErrorHandler(errorModel, err, "Please provide valid CVC.", "", "ValidateTicketPaymentRequest")
		}
	}

	
	//ticket validation
	if(len(request.Tickets)<1){
		err := errors.New("Please provide ticket detail.");
		if(err != nil){
			errorModel = ErrorHandler(errorModel, err, "Please provide ticket detail.", "", "ValidateTicketPaymentRequest")
		}
	}else{
		for i, ticket := range request.Tickets {
			fmt.Println(i);
			
			if(uuid.FromStringOrNil(ticket.TicketTypeId) == uuid.Nil){
				err := errors.New("Please provide valid TicketTypeId.");
				if(err != nil){
					errorModel = ErrorHandler(errorModel, err, "Please provide valid TicketTypeId.", "", "ValidateTicketPaymentRequest")
				}
			}
			if(ticket.Quantity <1){
				err := errors.New("Please select at least 1 ticket against each selected ticket type.");
				if(err != nil){
					errorModel = ErrorHandler(errorModel, err, "Please select at least 1 ticket against each selected ticket type.", "", "ValidateTicketPaymentRequest")
				}
			}
		}
	}

	  return errorModel
  }

