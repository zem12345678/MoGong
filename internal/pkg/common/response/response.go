package response

import (
	"mogong/internal/pkg/common/exceptions"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func JsonResponse(ctx *gin.Context, httpCode, code int, err, data interface{}) {
	ctx.JSON(httpCode, Response{
		Code: code,
		Msg:  exceptions.GetMsg(code),
		Data: data,
	})
}
