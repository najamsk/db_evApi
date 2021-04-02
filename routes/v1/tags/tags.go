package tags

import (
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/gin-gonic/gin"
	"github.com/najamsk/eventvisor/eventvisor.api/services"
	"github.com/satori/go.uuid"
	//"fmt"
)

var apiConfig *config.Config
var apiRouter *gin.Engine
//var roleService services.Role
//var userService services.User
//var imageService services.Image
var tagService services.Tag

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration

	apiRouter = router
	api := router.Group("api/v1/tags")
	{
		
		api.POST("/all", utils.BasicAuth(configuration), getUserTags)
		api.POST("/update", utils.BasicAuth(configuration), updateTags)
	}
	apiv2 := router.Group("api/v2/tags")
	{
		
		apiv2.POST("/all", utils.BasicAuthV2(configuration), getUserTags)
		apiv2.POST("/update", utils.BasicAuthV2(configuration), updateTags)
	}

}

func getTags(c *gin.Context) {
	var taglist []models.Tag
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; ClientID uuid.UUID}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var tagsdb, err = tagService.GetAll();

		if(err != nil){
			taglist = nil;
			loggy.Logger.Info().Msg("getTags:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Tags not found.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		}else{
			taglist = tagsdb
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"tags": taglist,
		},
	)
}

func updateTags(c *gin.Context) {
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	type RequestModel struct {UserID uuid.UUID; TagIds []string}
	var req RequestModel;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var err = tagService.UpdateUserTags(req.UserID, req.TagIds);

		if(err != nil){
			loggy.Logger.Info().Msg("getTags:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Error while updating user tags. Please try again.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		}else{success_message="Tag(s) updated successfully."}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message,
		},
	)
}

func getUserTags(c *gin.Context) {
	var err = error(nil);
	var taglist []*viewmodels.TagVM
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var tagsdb []*viewmodels.TagVM;
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; ClientID uuid.UUID}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
	//req.ConferenceID, err := uuid.FromString(u)
	if(req.ConferenceID == uuid.Nil){
		req.ConferenceID, err = uuid.FromString("b10f4c67-930e-4e76-93fd-ddd52c532b79");
	}
		 tagsdb, err = tagService.GetAllTagsWithSelected(req.UserID, req.ConferenceID);

		if(err != nil){
			taglist = nil;
			loggy.Logger.Info().Msg("getTags:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Tags not found.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		}else{
			taglist = tagsdb
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"tags": taglist,
		},
	)
}
