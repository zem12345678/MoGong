package models

import (
	"encoding/json"
	"mogong/global"
)

type Settings struct {
	global.MG_MODEL
	Classify int
	Content  json.RawMessage
}

func (s Settings) TableName() string {
	return "sys_settings"
}
