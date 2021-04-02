package favorites

import (
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	//"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/gin-gonic/gin"
	"github.com/najamsk/eventvisor/eventvisor.api/services"

	//"github.com/satori/go.uuid"
)

var apiConfig *config.Config
var apiRouter *gin.Engine
//var roleService services.Role
//var userService services.User
//var imageService services.Image
var favoriteService services.Favorite
var sessionService services.Session

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration

	apiRouter = router
	api := router.Group("api/v1/favorites")
	{
		api.POST("/update", utils.BasicAuth(configuration), updateFavorotes)
		api.POST("/get", utils.BasicAuth(configuration), getFavorotes)	
	}
	apiv2 := router.Group("api/v2/favorites")
	{
		apiv2.POST("/update", utils.BasicAuthV2(configuration), updateFavorotes)
		apiv2.POST("/get", utils.BasicAuthV2(configuration), getFavorotes)	
	}
}
func updateFavorotes(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var errorvm *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID; SpeakerID uuid.UUID; Favorite bool;}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
	if(req.Favorite == false){
		var err = favoriteService.RemoveFavorite(req.ConferenceID, req.UserID, req.EntityID, req.EntityType);
		
		if(err != nil){
			loggy.Logger.Info().Msg("updateFavorotes:Error: "+err.Error())
			if(errorvm == nil){
				errorvm = new(viewmodels.Error)
				}
				errorvm.DisplayErrors = append(errorvm.DisplayErrors, "Not updated. Please try again.")
				errorvm.InnerErrors = append(errorvm.InnerErrors, err.Error())
			//status_code = http.StatusBadRequest
			errorvm.ApiStatusCode = http.StatusBadRequest
			
		}

	}else if(req.Favorite == true){

		var capacityReached = false
		var err = error(nil)

		if(req.EntityType == "sessions"){
				capacityReached, err = sessionService.CapacityReached(req.EntityID, req.UserID);}

		if(err == nil){
			if(capacityReached == false){
				err = favoriteService.AddFavorite(req.ConferenceID, req.UserID, req.EntityID, req.EntityType)
			}else{
				if(errorvm == nil){
					errorvm = new(viewmodels.Error)
					}
					errorvm.DisplayErrors = append(errorvm.DisplayErrors, "Session is full now you are on wish list.")
					errorvm.InnerErrors = append(errorvm.InnerErrors, "Session is full now you are on wish list.")
					errorvm.ApiStatusCode = http.StatusBadRequest
			}
		}
		
		if(err != nil){
			loggy.Logger.Info().Msg("updateFavorotes:Error: "+err.Error())
			if(errorvm == nil){
				errorvm = new(viewmodels.Error)
				}
				errorvm.DisplayErrors = append(errorvm.DisplayErrors, "Not updated. Please try again.")
				errorvm.InnerErrors = append(errorvm.InnerErrors, err.Error())
				errorvm.ApiStatusCode = http.StatusBadRequest
		}	
	}	

	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":errorvm, 
			"success_message": success_message, 
		},
	)
}

func getFavorotes(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var attendees []viewmodels.AttendeeVM
	var sessions []viewmodels.SessionVM
	//var empty_message_session string
	//var empty_message_attendees string
	//var empty_message_speaker string
	
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID; SpeakerID uuid.UUID; Favorite bool;}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
		
	var speakers, err  = favoriteService.GetFavoriteSpeakers(req.ConferenceID, req.UserID);
	attendees, err  = favoriteService.GetFavoriteAttendees(req.ConferenceID, req.UserID);
	sessions, err   = favoriteService.GetFavoriteSessions(req.ConferenceID, req.UserID);
		
		if(err != nil){
			loggy.Logger.Info().Msg("updateFavorotes:Error: "+err.Error())
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "Error whiel getting data.")
			error.InnerErrors = append(error.InnerErrors, err.Error())
			//status_code = http.StatusBadRequest
			error.ApiStatusCode = http.StatusBadRequest
			
		}
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"speakers":speakers,
			"attendees":attendees,
			"sessions":sessions,
			"empty_message_session": "You do not have any session in your agenda.",
			"empty_message_attendees":"You do not have any attendee(s) in your favorite list.",
			"empty_message_speaker":"You do not have any speaker(s) in your favorite list.",
		},
	)
}


