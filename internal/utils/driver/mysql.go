package driver

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMysql(host, username, password, dbname, port string, logMode logger.LogLevel,
	maxIdle, maxOpen int, maxLifetime time.Duration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 最大閒置連線數
	sqlDB.SetMaxIdleConns(maxIdle)
	// 最大連線數
	sqlDB.SetMaxOpenConns(maxOpen)
	// 每條連線的存活時間
	sqlDB.SetConnMaxLifetime(maxLifetime)

	return db, nil
}
