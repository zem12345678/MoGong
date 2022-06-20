package models

import "mogong/global"

type SysPost struct {
	global.MG_MODEL
	Name        string `json:"name" gorm:"comment:岗位名称"`
	Code        string `gorm:"size:128;" json:"postCode"` //岗位代码
	Sort        int    `gorm:"size:4;" json:"sort"`       //岗位排序
	Status      int    `gorm:"size:4;" json:"status"`     //状态
	Description string `gorm:"size:255;" json:"remark"`   //描述
	DataScope   string `gorm:"-" json:"dataScope"`
	Params      string `gorm:"-" json:"params"`
	global.OperateBy
}

func (SysPost) TableName() string {
	return "sys_post"
}
