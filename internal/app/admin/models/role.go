package models

import "mogong/global"

type SysRole struct {
	global.MG_MODEL
	Name          string    `json:"name" gorm:"comment:角色名"`        // 角色名
	ParentID      int64     `json:"parent_id" gorm:"comment:父角色ID"` // 父角色ID
	SuperAdmin    bool      `json:"super_admin" gorm:"comment:"`
	DataScope     string    `json:"dataScope" gorm:"size:128;"`
	Children      []SysRole `json:"children" gorm:"-"`
	MenuIds       []int     `json:"menuIds" gorm:"-"`
	DeptIds       []int     `json:"deptIds" gorm:"-"`
	SysDept       []SysDept `json:"sysDept" gorm:"many2many:sys_role_dept;foreignKey:ID;joinForeignKey:id;references:ID;joinReferences:id;"`
	Menus         []SysMenu `json:"sysMenu" gorm:"many2many:sys_role_menu;foreignKey:ID;joinForeignKey:id;references:ID;joinReferences:id;"`
	DefaultRouter string    `json:"default_router" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认dashboard)
	global.OperateBy
}
