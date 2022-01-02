package database

import (
	"fmt"
	"todo-apis-go/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Declare the variable for the database
var DB *gorm.DB

// ConnectDB connect to db
func ConnectDB() {
	var err error
	var dsn string

	// if !fiber.IsChild() {
	// fmt.Println("I'm the parent process")
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?&parseTime=True&charset=utf8&loc=Local", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), 3306)
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         100,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		DontSupportForShareClause: true,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.Config("MYSQL_DBNAME")))
	DB.Exec(fmt.Sprintf("USE %s", config.Config("MYSQL_DBNAME")))
	DB.Exec("CREATE TABLE IF NOT EXISTS `todos` (`id` int AUTO_INCREMENT,`title` varchar(50),`activity_group_id` int unsigned,`is_active` boolean,`priority` varchar(10),`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,PRIMARY KEY (`id`),INDEX idx_todos_activity_group_id (`activity_group_id`),INDEX idx_todos_deleted_at (`deleted_at`))")
	DB.Exec("CREATE TABLE IF NOT EXISTS `activities` (`id` int AUTO_INCREMENT,`email` varchar(50),`title` varchar(50),`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,PRIMARY KEY (`id`),INDEX idx_activities_deleted_at (`deleted_at`))")

	sqlDB, _ := DB.DB()
	sqlDB.Close()

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&parseTime=True&charset=utf8&loc=Local", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), 3306, config.Config("MYSQL_DBNAME"))
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         100,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		DontSupportForShareClause: true,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic("failed to connect database")
	}

}
