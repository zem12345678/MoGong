package models

type CasbinRule struct {
	PType  string `json:"p_type" gorm:"column:p_type"`
	Role   string `json:"role_name" gorm:"column:v0"`
	Path   string `json:"path" gorm:"column:v1"`
	Method string `json:"method" gorm:"column:v2"`
	V3     string `json:"v3" gorm:"size:100;"`
	V4     string `json:"v4" gorm:"size:100;"`
	V5     string `json:"v5" gorm:"size:100;"`
}

func (c CasbinRule) TableName() string {
	return "casbin_rule"
}
