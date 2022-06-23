package controllers

import (
	"mogong/internal/app/admin/services"
	"mogong/internal/pkg/common/exceptions"
	"mogong/internal/pkg/common/models/user/request"
	"mogong/internal/pkg/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	logger  *zap.Logger
	service services.SysUserService
}

func NewUserController(logger *zap.Logger, s services.SysUserService) *UserController {
	return &UserController{
		logger:  logger.With(zap.String("type", "UserController")),
		service: s,
	}
}

func (u *UserController) Login(c *gin.Context) {
	req := new(request.LoginRequest)
	if err := c.ShouldBindJSON(req); err != nil {

	}
	res, err := u.service.Login(req)
	if err != nil {

	}
	response.JsonResponse(c, http.StatusOK, exceptions.Success, nil, res)
	return

}
