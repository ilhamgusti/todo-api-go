package database

import (
	"database/sql"
	"fmt"
	"time"
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

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?&parseTime=True&charset=utf8&loc=Local", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), 3306)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + config.Config("MYSQL_DBNAME"))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + config.Config("MYSQL_DBNAME"))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `todos` (`id` int AUTO_INCREMENT,`title` varchar(50),`activity_group_id` int unsigned,`is_active` boolean,`priority` varchar(10),`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,PRIMARY KEY (`id`),INDEX idx_todos_activity_group_id (`activity_group_id`),INDEX idx_todos_deleted_at (`deleted_at`))")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `activities` (`id` int AUTO_INCREMENT,`email` varchar(50),`title` varchar(50),`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,PRIMARY KEY (`id`),INDEX idx_activities_deleted_at (`deleted_at`))")
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&parseTime=True&charset=utf8&loc=Local", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), 3306, config.Config("MYSQL_DBNAME"))
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	sqlDB, _ := DB.DB()
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetMaxOpenConns(300)

	if err != nil {
		panic("failed to connect database")
	}

}
