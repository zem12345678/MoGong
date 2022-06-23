package constants

import (
	"mogong/internal/pkg/common/exceptions"
)

var (
	NotFoundUserErr = exceptions.NewBaseError("用户不存在")
	UserPasswordErr = exceptions.NewBaseError("用户不存在或者密码错误")
	// AccessTokenErr 生成签名错误
	AccessTokenErr = exceptions.NewBaseError("生成签名错误")
)
