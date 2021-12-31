package database

import (
	"fmt"
	"todo-apis-go/config"
	"todo-apis-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Declare the variable for the database
var DB *gorm.DB

// ConnectDB connect to db
func ConnectDB() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?&parseTime=True", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), 3306)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic("failed to connect database")
	}

	DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.Config("MYSQL_DBNAME")))
	DB.Exec(fmt.Sprintf("USE %s", config.Config("MYSQL_DBNAME")))
	errs := DB.AutoMigrate(&models.Todo{}, &models.Activity{})

	if errs != nil {
		fmt.Sprintln(errs)
	}

	sqlDB, err := DB.DB()
	sqlDB.Close()

	if err != nil {
		panic("failed to connect database")
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&parseTime=True", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), 3306, config.Config("MYSQL_DBNAME"))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic("failed to connect database")
	}
}
