package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
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
var clientService services.Client

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration

	apiRouter = router
	api := router.Group("api/v1/clients")
	{
		
		api.POST("/get", utils.BasicAuth(configuration), getClient)
		api.GET("/test", utils.BasicAuth(configuration), testGet)
		api.POST("/contact", utils.BasicAuth(configuration), getContact)
	}
	apiv2 := router.Group("api/v2/clients")
	{
		
		apiv2.POST("/get", utils.BasicAuthV2(configuration), getClient)
		apiv2.GET("/test", utils.BasicAuthV2(configuration), testGet)
		apiv2.POST("/contact", utils.BasicAuthV2(configuration), getContact)
	}

}
func printStr(v string) {
	fmt.Println(v)
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Println(string(vBytes))
}
func getContact(c *gin.Context) {

	c.JSON(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"phone": apiConfig.Items.Contact.Phone,
			"email": apiConfig.Items.Contact.Email,
			"web": apiConfig.Items.Contact.Web,
			"Phone2": apiConfig.Items.Contact.Phone2,
			"WebDisplay" : apiConfig.Items.Contact.WebDisplay,
		},
	)
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

func getClient(c *gin.Context) {
	var clientinfo *models.Client
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; ClientID uuid.UUID}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var clientdb, err = clientService.Get(req.ClientID);

		if(err != nil){
			clientinfo = nil;
			loggy.Logger.Info().Msg("getProfile:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Client does not exist.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		}else{
			clientinfo = &clientdb
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"client": clientinfo,
		},
	)
}
