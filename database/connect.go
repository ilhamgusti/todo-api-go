package database

import (
	"database/sql"
	"fmt"
	"strconv"
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
	port, err := strconv.Atoi(config.Config("MYSQL_PORT"))
	if err != nil {
		panic("Port Not Same")
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?&parseTime=True&charset=utf8&loc=Local", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), port)
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `todos` (`id` varchar(50),`title` varchar(50),`activity_group_id` int unsigned,`is_active` boolean,`priority` varchar(10)) ENGINE=MyISAM")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `activities` (`id` varchar(50),`email` varchar(50),`title` varchar(50)) ENGINE=MyISAM")
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&parseTime=True&charset=utf8&loc=Local", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), port, config.Config("MYSQL_DBNAME"))
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DontSupportRenameIndex:    true,
		DisableDatetimePrecision:  true,
		DontSupportRenameColumn:   true,
		DontSupportForShareClause: true,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		DisableAutomaticPing:   true,
	})

	if err != nil {
		panic("failed to connect database")
	}

}
