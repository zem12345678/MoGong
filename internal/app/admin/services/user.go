package services

import (
	"mogong/internal/app/admin/constants"
	"mogong/internal/app/admin/repositories"
	"mogong/internal/pkg/common/models/user/request"
	"mogong/internal/pkg/common/models/user/response"
	"mogong/internal/pkg/tools/jwt"
	"mogong/internal/pkg/tools/make_password"
	"time"

	"go.uber.org/zap"
)

type SysUserService interface {
	Login(req *request.LoginRequest) (*response.LoginResponse, error)
}

type DefaultSysUserService struct {
	logger *zap.Logger
	repo   repositories.SysUserRepo
}

func NewUserService(logger *zap.Logger, repo repositories.SysUserRepo) SysUserService {
	return &DefaultSysUserService{
		logger: logger.With(zap.String("type", "DefaultSysUserService")),
		repo:   repo,
	}
}

func (d *DefaultSysUserService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	user, err := d.repo.GetUserByUserName(req.UserName)
	if err != nil {
		return nil, constants.NotFoundUserErr
	}
	if !make_password.CheckPassword(req.Password, user.Password) {
		return nil, constants.UserPasswordErr
	}
	accessToken, expired, err := jwt.CreateToken(req.UserName, user.ID)
	if err != nil {
		return nil, constants.AccessTokenErr
	}
	return &response.LoginResponse{
		AccessToken: accessToken,
		ExpireAt:    expired,
		TimeStamp:   time.Now().Unix(),
	}, nil
}
