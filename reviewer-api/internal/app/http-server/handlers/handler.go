package handlers

import (
	"reviewer-api/internal/app/http-server/handlers/pk"
	"reviewer-api/internal/app/http-server/handlers/team"
	"reviewer-api/internal/app/http-server/handlers/user"
	"reviewer-api/internal/app/repository/postgres"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Team *team.TeamHandler
	User *user.UserHandler
	PR   *pk.PKHandler
}

func NewHandler(pg *postgres.Postgres) (*Handlers, error) {
	return &Handlers{
		Team: team.NewTeamHandler(pg),
		User: user.NewUserHandler(pg),
		PR:   pk.NewPKHandler(pg),
	}, nil
}

func (h *Handlers) Register(r *gin.Engine) {
	r.GET("/team/get", h.Team.GetTeam)
	r.POST("/team/add", h.Team.AddTeam)

	r.POST("/users/setIsActive", h.User.UpdateUserActivity)
	r.GET("/users/getReview", h.User.GetUserReview)

	r.POST("/pullRequest/create", h.PR.CreateNewPullRequest)
	r.POST("/pullRequest/reassign", h.PR.ReassignPullRequest)
	r.POST("/pullRequest/merge", h.PR.MergedPR)
}
