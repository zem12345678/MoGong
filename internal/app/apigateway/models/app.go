package models

import "mogong/global"

type App struct {
	global.MG_MODEL
	AppID    string `json:"app_id" gorm:"column:app_id" description:"租户id"`
	Name     string `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret   string `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS string `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd      int64  `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps      int64  `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
}

func (App) TableName() string {
	return "gateway_app"
}
