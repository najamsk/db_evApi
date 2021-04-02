package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	//"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin"
	"github.com/najamsk/eventvisor/eventvisor.api/services"
	//"os"
	//"io"
)

var apiConfig *config.Config
var apiRouter *gin.Engine
var sessionService services.Session

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration

	apiRouter = router
	api := router.Group("api/v1/session")
	{
		
		api.POST("/get", utils.BasicAuth(configuration), getSession)
		api.POST("/rating/update", utils.BasicAuth(configuration), updateSessionRating)
	}
	apiv2 := router.Group("api/v2/session")
	{
		
		apiv2.POST("/get", utils.BasicAuthV2(configuration), getSession)
		apiv2.POST("/rating/update", utils.BasicAuthV2(configuration), updateSessionRating)
	}

}
func printStr(v string) {
	fmt.Println(v)
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Println(string(vBytes))
}

func getSession(c *gin.Context) {
	//var clientinfo *models.Client
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	type RequestModel struct {ConferenceID uuid.UUID; UserID uuid.UUID; SessionID uuid.UUID}
	var req RequestModel;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var sessionDB, err = sessionService.GetSession(req.SessionID, req.ConferenceID, req.UserID);

		if(err != nil){
			//clientinfo = nil;
			loggy.Logger.Info().Msg("getSession:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Session does not exist.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		}else{
			//clientinfo = &clientdb
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"session": sessionDB,
		},
	)
}

func updateSessionRating(c *gin.Context) {
	//var clientinfo *models.Client
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	type RequestModel struct {ConferenceId uuid.UUID; SessionID uuid.UUID; UserID uuid.UUID; Rating float64}
	var req RequestModel;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var err = sessionService.UpdateRating(req.UserID, req.SessionID, req.Rating);
		var avgRating = sessionService.GetSessionAvgRating(req.SessionID, req.UserID);

		if(err != nil){
			loggy.Logger.Info().Msg("updateSessionRating:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Session rating not updated.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		error.ApiStatusCode = http.StatusBadRequest
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"avg_rating": avgRating,
		},
	)
}
