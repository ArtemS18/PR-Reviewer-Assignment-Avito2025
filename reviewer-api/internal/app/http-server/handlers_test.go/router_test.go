package handlers_test

import (
	"reviewer-api/internal/app/http-server/handlers"
	"reviewer-api/internal/app/http-server/handlers/pk"
	"reviewer-api/internal/app/http-server/handlers/team"
	"reviewer-api/internal/app/http-server/handlers/user"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	h := &handlers.Handlers{
		Team: team.NewTeamHandler(&mockTeamRepo{}),
		User: user.NewUserHandler(&mockUserRepo{}),
		PR:   pk.NewPKHandler(&mockPKRepo{}),
	}
	h.Register(r)

	return r
}
