package models

import (
	"errors"

	"gorm.io/gorm"
)

type DataPermission struct {
	DataScope string `json:"data_scope"`
	UserId    int    `json:"user_id"`
	DeptId    int    `json:"dept_id"`
	RoleId    int    `json:"role_id"`
}

func (d *DataPermission) GetDataScope(tableName string, db *gorm.DB) (*gorm.DB, error) {
	user := new(SysUser)
	role := new(SysRole)
	err := db.Find(user, d.UserId).Error
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	err = db.Find(role, user.RoleId).Error
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	if role.DataScope == "2" {
		db = db.Where(tableName+".create_by in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
	}
	if role.DataScope == "3" {
		db = db.Where(tableName+".create_by in (SELECT user_id from sys_user where dept_id = ? )", user.DeptId)
	}
	if role.DataScope == "4" {
		db = db.Where(tableName+".create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+pkg.IntToString(user.DeptId)+"%")
	}
	if role.DataScope == "5" || role.DataScope == "" {
		db = db.Where(tableName+".create_by = ?", d.UserId)
	}

	return db, nil
}
