package routes

import (
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	users "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/users"
	clients "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/clients"
	conferences "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/conferences"
	speakers "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/speakers"
	images "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/images"
	tags "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/tags"
	favorites "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/favorites"
	sessions "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/sessions"
	tickets "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/tickets"
	payments "github.com/najamsk/eventvisor/eventvisor.api/routes/v1/payments"

	"github.com/gin-gonic/gin"
	"fmt"
	"time"
	"strconv"
	"math/rand"
)

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		now := time.Now()
		sec := now.Unix()
		fmt.Println("time1:",sec)
		fmt.Println("time:",strconv.FormatInt(sec, 10))

		rand.Seed(int64(time.Now().Nanosecond()))
		myrand := rand.Intn(999999 - 100000) + 100000
		fmt.Println(myrand)
		// Call the HTML method of the Context to render a template
		c.JSON(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "Home Page",
			},
		)

	})

	users.Routes(configuration, router)
	clients.Routes(configuration, router)
	conferences.Routes(configuration, router)
	speakers.Routes(configuration, router)
	images.Routes(configuration, router)
	tags.Routes(configuration, router)
	favorites.Routes(configuration, router)
	sessions.Routes(configuration, router)
	tickets.Routes(configuration, router)
	payments.Routes(configuration, router)


}
