package utils

import (
	 "github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	"github.com/gin-gonic/gin"
	"fmt"
	"strings"
	"encoding/base64"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)
// AuthUserKey is the cookie name for user credential in basic auth.
const AuthUserKey = "user"
//var userService User
//BasicAuth ead these from config. and setup a func in utils to export this func
func BasicAuth(configuration *config.Config,) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		configuration.Items.BasicAuth.Username: configuration.Items.BasicAuth.Password,
		//"foo": "bar",
	})
}

func BasicAuthV2(configuration *config.Config) gin.HandlerFunc {
	
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		  }

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		fmt.Println("pair", pair);

		if len(pair) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		  }

		var email string = pair[0];
		var password string = pair[1];
	
		var userdb models.User
		db := GetDb();
		fmt.Println("email:",email);
		fmt.Println("password:",password);
		var err = db.Preload("Roles").Where("lower(email) = lower(?) and is_active = ?", email, true).First(&userdb).Error
		if err != nil{
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userdb.Password), []byte(password))
		if(err != nil){
			c.AbortWithStatus(http.StatusUnauthorized)
			return			
		}

		var validRole bool = IsInRole(userdb.Roles, []string{"ClientApi"});
		fmt.Println("validRole:",validRole);
		if validRole != true{
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(AuthUserKey, userdb)
	return;
}
}
