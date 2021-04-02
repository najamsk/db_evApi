package speakers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	//"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	//"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin"
	"github.com/najamsk/eventvisor/eventvisor.api/services"
	//"os"
	//"io"
)

var apiConfig *config.Config
var apiRouter *gin.Engine
//var roleService services.Role
//var userService services.User
//var imageService services.Image
var speakerService services.Speaker

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration

	apiRouter = router
	api := router.Group("api/v1/speakers")
	{
		
		api.POST("/getall", utils.BasicAuth(configuration), getConferenceSpeakers)
		api.POST("/profile", utils.BasicAuth(configuration), getSpeakerProfile)
		api.GET("/test", utils.BasicAuth(configuration), testGet)
	}

	apiv2 := router.Group("api/v2/speakers")
	{
		
		apiv2.POST("/getall", utils.BasicAuthV2(configuration), getConferenceSpeakers)
		apiv2.POST("/profile", utils.BasicAuthV2(configuration), getSpeakerProfile)
		apiv2.GET("/test", utils.BasicAuthV2(configuration), testGet)
	}

}
func printStr(v string) {
	fmt.Println(v)
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Println(string(vBytes))
}
func testGet(c *gin.Context) {

	c.JSON(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "Test Get",
		},
	)
}

func getConferenceSpeakers(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID}
	var req viewmodels.Request;
	var pageSize int = 50
	//var offSet int = 0
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)

	// if(req.OffSet<1){
	// 	req.OffSet = 0;
	// }
	
		if(req.PageNo<1){
			req.OffSet = -1
		}
	var speakers, err = speakerService.GetSpeakersByConference(req.ConferenceID, req.UserID, pageSize, req.OffSet, req.SearchText);
	fmt.Println(err)
		
	defer loggy.CloseLog()

	// if(len(speakers)<pageSize){
	// 	req.PageNo = -1
	// }else if(req.PageNo>0){
	// 	req.PageNo = req.PageNo +1
	// }

	if(len(speakers) == pageSize && req.PageNo>0){
		req.PageNo = req.PageNo +1
	}
	// if(len(speakers)<pageSize){
	// 	req.OffSet = (req.PageNo-1) * pageSize -(pageSize - len(speakers))
	// }else{
	// 	req.OffSet = pageSize*(req.PageNo-1)
	// }

	req.OffSet = req.OffSet + len(speakers);
			
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"speakers": speakers,
			"next_page_no": req.PageNo,
			"speakers_count": len(speakers),
			"off_set":req.OffSet,
		},
	)
}

func getSpeakerProfile(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID; SpeakerID uuid.UUID;}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
	var profile, err = speakerService.GetSpeakerProfile(req.ConferenceID, req.UserID, req.SpeakerID);
	fmt.Println(err)

	var sessions, err_session = speakerService.GetSpeakerSessions(req.ConferenceID, req.UserID, req.SpeakerID)
	fmt.Println(err_session)		
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"profile": profile,
			"sessions": sessions,
		},
	)
}

// func updateFavoroteSpeaker(c *gin.Context) {
// 	//var conferenceinfo *models.Conference
// 	var error *viewmodels.Error
// 	var success_message string
// 	var status_code = http.StatusOK
	
// 	var loggy = utils.FLogger{}
// 	loggy.OpenLog()
// 	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID; SpeakerID uuid.UUID; Favorite bool;}
// 	var req viewmodels.Request;
	
// 	//c.BindJSON(&userinfo)
// 	c.BindJSON(&req)
	
// 	if(req.Favorite == false){
// 		var err = speakerService.RemoveFavoriteSpeaker(req.ConferenceID, req.UserID, req.SpeakerID);
		
// 		if(err != nil){
// 			loggy.Logger.Info().Msg("RegisterUser:Error: "+err.Error())
// 			if(error == nil){
// 				error = new(viewmodels.Error)
// 				}
// 			error.DisplayErrors = append(error.DisplayErrors, "Not updated. Please try again.")
// 			error.InnerErrors = append(error.InnerErrors, err.Error())
// 			//status_code = http.StatusBadRequest
// 			error.ApiStatusCode = http.StatusBadRequest
			
// 		}

// 	}else if(req.Favorite == true){
// 		var err = speakerService.AddFavoriteSpeaker(req.ConferenceID, req.UserID, req.SpeakerID)
// 		if(err != nil){
// 			loggy.Logger.Info().Msg("RegisterUser:Error: "+err.Error())
// 			if(error == nil){
// 				error = new(viewmodels.Error)
// 				}
// 			error.DisplayErrors = append(error.DisplayErrors, "Not updated. Please try again.")
// 			error.InnerErrors = append(error.InnerErrors, err.Error())
// 			//status_code = http.StatusBadRequest
// 			error.ApiStatusCode = http.StatusBadRequest
			
// 		}	
// 	}	

// 	defer loggy.CloseLog()
	
// 	c.JSON(
// 		status_code,
// 		gin.H{
// 			"error":error, 
// 			"success_message": success_message, 
// 		},
// 	)
// }





