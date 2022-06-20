package models

import "mogong/global"

type SysApi struct {
	global.MG_MODEL
	Path        string `json:"path" gorm:"comment:api路径"`            // api路径
	Description string `json:"description" gorm:"comment:api中文描述"`   // api中文描述
	ApiGroup    string `json:"api_group" gorm:"comment:api组"`        // api组
	Method      string `json:"method" gorm:"default:GET;comment:方法"` // 方法:创建POST|查看GET(默认)|更新PUT|删除DELETE
	global.OperateBy
}

func (SysApi) TableName() string {
	return "sys_apis"
}
