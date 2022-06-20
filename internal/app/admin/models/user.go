package models

import (
	"mogong/global"

	uuid "github.com/satori/go.uuid"
)

type SysUser struct {
	global.MG_MODEL
	UUID        uuid.UUID `json:"uuid" gorm:"comment:用户UUID"`                                                            // 用户UUID
	UserName    string    `json:"user_name" gorm:"comment:用户登录名"`                                                        // 用户登录名
	Password    string    `json:"-"  gorm:"comment:用户登录密码"`                                                              // 用户登录密码
	NickName    string    `json:"nick_name" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	SideMode    string    `json:"side_node" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	HeaderImg   string    `json:"header_img" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	BaseColor   string    `json:"base_color" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	ActiveColor string    `json:"active_color" gorm:"default:#1890ff;comment:活跃颜色"`
	Sex         string    `json:"sex" gorm:"size:8;comment:性别"`
	Email       string    `json:"email" gorm:"size:128;comment:邮箱"`
	Phone       string    `json:"phone"  gorm:"comment:用户手机号"`
	Description string    `json:"description" gorm:"comment:api中文描述"` // 描述
	RoleId      int       `json:"roleId" gorm:"size:20;comment:角色ID"`
	DeptId      int       `json:"dept_id" gorm:"size:20;comment:部门"`
	PostId      int       `json:"post_id" gorm:"size:20;comment:岗位"`
	Status      string    `json:"status" gorm:"size:4;comment:状态"`
	DeptIds     []int     `json:"dept_ids" gorm:"-"`
	PostIds     []int     `json:"post_ids" gorm:"-"`
	RoleIds     []int     `json:"role_ids" gorm:"-"`
	Dept        *SysDept  `json:"dept"`
	global.OperateBy
}

func (SysUser) TableName() string {
	return "sys_user"
}
