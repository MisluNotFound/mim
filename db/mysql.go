package db

import (
	"gorm.io/gorm"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("init database failed: ", zap.Error(err))
		return err
	}
	DB = db
	return nil
}
