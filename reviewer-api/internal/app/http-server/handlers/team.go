package handlers

import (
	"net/http"
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	pkg "reviewer-api/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type TeamRepository interface {
	GetTeam(team_name string) (ds.Team, error)
	AddTeam(teamData dto.TeamDTO) (ds.Team, error)
}

type TeamHandler struct {
	repo TeamRepository
}

func (h *TeamHandler) GetTeam(ctx *gin.Context) {
	teamName := ctx.Query("team_name")
	if teamName == "" {
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			pkg.NOT_FOUND,
		)
		return
	}
	team, err := h.repo.GetTeam(teamName)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	pkg.OkResponse(ctx, team)
}

func (h *TeamHandler) AddTeam(ctx *gin.Context) {
	var teamDTO dto.TeamDTO
	err := ctx.BindJSON(&teamDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			pkg.BAD_REQUEST,
		)
	}
	team, err := h.repo.AddTeam(teamDTO)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, team)
}
