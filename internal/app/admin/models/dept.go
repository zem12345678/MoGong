package models

import "mogong/global"

type SysDept struct {
	global.MG_MODEL
	ParentID  int       `json:"parent_Id" gorm:""`       //上级部门
	Path      string    `json:"path" gorm:"size:255;"`   //
	Name      string    `json:"name"  gorm:"size:128;"`  //部门名称
	Sort      int       `json:"sort" gorm:"size:4;"`     //排序
	Leader    string    `json:"leader" gorm:"size:128;"` //负责人
	Phone     string    `json:"phone" gorm:"size:11;"`   //手机
	Email     string    `json:"email" gorm:"size:64;"`   //邮箱
	Status    int       `json:"status" gorm:"size:4;"`   //状态
	DataScope string    `json:"dataScope" gorm:"-"`
	Params    string    `json:"params" gorm:"-"`
	Children  []SysDept `json:"children" gorm:"-"`
	global.OperateBy
}

func (SysDept) TableName() string {
	return "sys_dept"
}
