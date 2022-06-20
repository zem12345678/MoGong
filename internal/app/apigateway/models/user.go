package models

import (
	"mogong/global"
)

type User struct {
	global.MG_MODEL
	UserName string `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt     string `json:"salt" gorm:"column:salt" description:"盐"`
	Password string `json:"password" gorm:"column:password" description:"密码"`
}

func (User) TableName() string {
	return "gateway_user"
}
