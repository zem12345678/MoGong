package database

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func GetReadDB(DBArray []*gorm.DB, length int) *gorm.DB {
	if len(DBArray) <= 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	return DBArray[rand.Intn(length)]
}
