package utils

import (
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	//"github.com/gin-gonic/gin"
	//"github.com/satori/go.uuid"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/token"
	//"encoding/json"
	//"strconv"
	"net/http"
	"fmt"
)
//first send call detail to stripe. stripe send card token 
//then send payment/charge call to stripe and pass card token
func MakePaymentCreditCardStripe(stripeKey string, CardNumber string, ExpiryMonth string, ExpiryYear string, CVC string, Amount float64, Currency string, Description string) (*stripe.Charge, *viewmodels.Error) {

	var errorModel *viewmodels.Error
	var cardToken string;
	//var success_message = "";
	stripe.Key = stripeKey;
	var paymentRes *stripe.Charge;
	var err = error(nil);

	var loggy = FLogger{}
	loggy.OpenLog()
	defer loggy.CloseLog()

	//param card detail to get card token so that we will send that token to payment transaction
	params := &stripe.TokenParams{
		Card: &stripe.CardParams{
		  Number: stripe.String(CardNumber),
		  ExpMonth: stripe.String(ExpiryMonth),
		  ExpYear: stripe.String(ExpiryYear),
		  CVC: stripe.String(CVC),
		},
	  }
	  cardTokenRes, err := token.New(params)

	  //check if stripe send token or there is some error in response
	  if(err != nil){
		loggy.Logger.Info().Msg("makeTicketPayment:card token response Error: "+err.Error())
		if(errorModel == nil){
			errorModel = new(viewmodels.Error)
			}
			errorModel.DisplayErrors = append(errorModel.DisplayErrors, "Invalid card.")
			errorModel.InnerErrors = append(errorModel.InnerErrors, err.Error())
			errorModel.ApiStatusCode = http.StatusBadRequest
	  }else{
		cardToken = cardTokenRes.ID;
	  }

	  //check if no errror then go for pament/charge
	  if(errorModel == nil){
		args := &stripe.ChargeParams{
			Amount: stripe.Int64(int64(Amount*100)),
			Currency: stripe.String(Currency),
			Description: stripe.String(Description),
		  }
		  args.SetSource(cardToken)
		  //args.SetIdempotencyKey("vEVG02v9zSrZ7lAk")
		  paymentRes, err = charge.New(args);
		  
		  if(err != nil){
			loggy.Logger.Info().Msg("makeTicketPayment: payment response Error: "+err.Error())
			errorModel = HandleStripeErrors(err, errorModel);
		  }
	  }

	  return paymentRes, errorModel
  }

  func HandleStripeErrors(err error,  errorModel *viewmodels.Error) *viewmodels.Error{
	fmt.Println("errorModel:",errorModel);
	//getting stripe specific error form generic error 
	if stripeErr, ok := err.(*stripe.Error); ok {
		// The Code field will contain a basic identifier for the failure.
		switch stripeErr.Code {
		case stripe.ErrorCodeAmountTooSmall:
			if(errorModel == nil){
				errorModel = new(viewmodels.Error)
				}
				errorModel.DisplayErrors = append(errorModel.DisplayErrors, "Please enter valid amount.")
				errorModel.InnerErrors = append(errorModel.InnerErrors, err.Error())
				errorModel.ApiStatusCode = http.StatusBadRequest
			default:
				if(errorModel == nil){
					errorModel = new(viewmodels.Error)
					}
					errorModel.DisplayErrors = append(errorModel.DisplayErrors, "Payment Failed. Try again.")
					errorModel.InnerErrors = append(errorModel.InnerErrors, err.Error())
					errorModel.ApiStatusCode = http.StatusBadRequest
		}
	}
	fmt.Println("errorModel:",errorModel);
	return errorModel
  }

