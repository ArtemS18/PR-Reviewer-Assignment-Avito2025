package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResp struct {
	Error ErrorData `json:"error"`
}

func OkResponse(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusOK, obj)
}

var NOT_FOUND = ErrorResp{
	Error: ErrorData{
		Code:    "NOT_FOUND",
		Message: "resource not found",
	},
}
var BAD_REQUEST = ErrorResp{
	Error: ErrorData{
		Code:    "BAD_REQUEST",
		Message: "bad data in request",
	},
}

var TEAM_EXISTS = ErrorResp{
	Error: ErrorData{
		Code:    "TEAM_EXISTS",
		Message: "team_name already exists",
	},
}

var PR_EXISTS = ErrorResp{
	Error: ErrorData{
		Code:    "PR_EXISTS",
		Message: "PR id already exists",
	},
}
