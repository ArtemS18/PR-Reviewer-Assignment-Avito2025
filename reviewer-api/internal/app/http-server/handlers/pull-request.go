package handlers

import (
	"net/http"
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
	pkg "reviewer-api/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type PKRepository interface {
	CreatePullRequest(pkDTO dto.PullRequestCreateDTO) (ds.PullRequest, error)
}

type PKHandler struct {
	repo PKRepository
}

func (h *PKHandler) CreateNewPullRequest(ctx *gin.Context) {
	var pkDTO dto.PullRequestCreateDTO
	err := ctx.BindJSON(&pkDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			pkg.BAD_REQUEST,
		)
	}
	team, err := h.repo.CreatePullRequest(pkDTO)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				pkg.NOT_FOUND,
			)
		case repository.ErrAlreadyExists:
			ctx.AbortWithStatusJSON(
				http.StatusConflict,
				pkg.PR_EXISTS,
			)
		default:
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				pkg.BAD_REQUEST,
			)
		}
		return
	}
	ctx.JSON(http.StatusCreated, dto.ToPullRequestDTO(team))
}
