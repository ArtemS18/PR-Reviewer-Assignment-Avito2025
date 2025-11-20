package handlers

import (
	"reviewer-api/internal/app/repository/postgres"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Team *TeamHandler
	User *UserHandler
	PR   *PKHandler
}

func NewHandler(pg *postgres.Postgres) (*Handlers, error) {
	return &Handlers{
		Team: &TeamHandler{pg},
		User: &UserHandler{pg},
		PR:   &PKHandler{pg},
	}, nil
}

func (h *Handlers) Register(r *gin.Engine) {
	r.GET("/team/get", h.Team.GetTeam)
	r.POST("/team/add", h.Team.AddTeam)

	r.POST("/users/setIsActive", h.User.UpdateUserActivity)
	r.GET("/users/getReview", h.User.GetUserReview)

	r.POST("/pullRequest/create", h.PR.CreateNewPullRequest)
}
