package database

import (
	"fmt"

	"mo_fiber_1/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() {
	var err error
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		SlowThreshold:             time.Second,
	// 		LogLevel:                  logger.Silent,
	// 		IgnoreRecordNotFoundError: true,
	// 		Colorful:                  false,
	// 	},
	// )
	DB, err = gorm.Open(sqlite.Open("database/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect local database")
	}
	SqlDB, _ = DB.DB()
	fmt.Println("connection openned to database")

	DB.AutoMigrate(&model.Product{}, &model.User{})
	fmt.Println("Database models migrated")
}
