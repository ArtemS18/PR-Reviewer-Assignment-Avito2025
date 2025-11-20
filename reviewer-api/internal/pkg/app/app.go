package pkg

import (
	"fmt"
	"log"
	"reviewer-api/internal/app/config"
	http "reviewer-api/internal/app/http-server/handlers"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *http.TeamHandler
}

func NewApplication(c *config.Config, r *gin.Engine, h *http.TeamHandler) *Application {
	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

func (app *Application) RunApplication() {
	log.Println("Server start up")
	//app.Handler.RegisterHandlers(app.Router)
	address := fmt.Sprintf("%s:%d", app.Config.HTTPHost, app.Config.HTTPPort)
	if err := app.Router.Run(address); err != nil {
		log.Fatal(err)
	}
	log.Println("Server down")
}
