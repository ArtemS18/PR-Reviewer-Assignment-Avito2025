package handlers

import (
	"net/http"
	"reviewer-api/internal/app/ds"
	pkg "reviewer-api/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	SetUserFlag(user_id string, is_active bool) (ds.User, error)
	GetReview(user_id string) (ds.User, error)
}

type UserHandler struct {
	repo UserRepository
}

type UserUpdateSchema struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandler) UpdateUserActivity(ctx *gin.Context) {
	var userData UserUpdateSchema
	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			pkg.BAD_REQUEST,
		)
	}
	user, err := h.repo.SetUserFlag(userData.UserId, userData.IsActive)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	pkg.OkResponse(ctx, user)
}

func (h *UserHandler) GetUserReview(ctx *gin.Context) {
	userName := ctx.Query("user_id")
	if userName == "" {
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			pkg.NOT_FOUND,
		)
		return
	}
	user, err := h.repo.GetReview(userName)
	if err != nil {
		pkg.HandelError(ctx, err)
		return
	}
	pkg.OkResponse(ctx, user)
}
