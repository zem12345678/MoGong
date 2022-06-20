package global

import (
	"time"

	"gorm.io/gorm"
)

type MG_MODEL struct {
	ID          int64          `gorm:"primaryKey" json:"id"`                     // 主键ID
	CreatedTime time.Time      `gorm:"autoUpdateTime:nano" json:"created_time" ` // 创建时间
	UpdatedTime time.Time      `gorm:"autoUpdateTime:nano" json:"updated_time"`  // 更新时间
	DeletedTime gorm.DeletedAt `gorm:"index" json:"-"`                           // 删除时间
}

type OperateBy struct {
	CreateBy int `json:"createBy" gorm:"index;comment:创建者"`
	UpdateBy int `json:"updateBy" gorm:"index;comment:更新者"`
}
