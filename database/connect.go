package database

import (
	"fmt"
	"log"
	"strconv"
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
	p := config.Config("MYSQL_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("Idiot")
	}

	// Connection URL to connect to Postgres Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?&parseTime=True", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), port)
	dsnwithDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&parseTime=True", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), port, config.Config("MYSQL_DBNAME"))
	// Connect to the DB and initialize the DB variable
	DB, err = gorm.Open(mysql.Open(dsn))

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connected to DATABASE")
	DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.Config("MYSQL_DBNAME")))
	DB.Exec(fmt.Sprintf("USE %s", config.Config("MYSQL_DBNAME")))

	fmt.Printf("Success Create Database %s ->", config.Config("MYSQL_DBNAME"))
	//Migrate the DB
	errs := DB.AutoMigrate(&models.Todo{}, &models.Activity{})
	if errs != nil {
		fmt.Sprintln(errs)
	}
	fmt.Println("Database Migrated!")
	fmt.Println("CLOSE CONNECTION DB")
	sqlDB, err := DB.DB()
	sqlDB.Close()

	fmt.Println("RECREATE NEW ONE...")
	DB, err = gorm.Open(mysql.Open(dsnwithDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Printf("Connected to EXISTING DATABASE: %s", config.Config("MYSQL_DBNAME"))
}
