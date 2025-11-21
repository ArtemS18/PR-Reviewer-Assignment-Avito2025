package pkg

import (
	"net/http"
	"reviewer-api/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func HandelError(ctx *gin.Context, err error) {
	switch err {
	case repository.ErrNotFound:
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			NOT_FOUND,
		)
	case repository.ErrTeamAlreadyExists:
		ctx.AbortWithStatusJSON(
			http.StatusConflict,
			TEAM_EXISTS,
		)
	case repository.ErrPRAlreadyExists:
		ctx.AbortWithStatusJSON(
			http.StatusConflict,
			PR_EXISTS,
		)
	default:
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			UNEXPECT,
		)
	}
}
