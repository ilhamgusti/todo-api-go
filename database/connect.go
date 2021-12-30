package database

import (
	"fmt"
	"log"
	"strconv"
	"todo-apis-go/config"
	"todo-apis-go/models"

	"github.com/gofiber/fiber/v2"
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
		log.Println("This is Port 'Number' not 'Kata Kata' Mutiara!")
	}

	if !fiber.IsChild() {
		fmt.Println("Creating Database on Parent Process...")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?&parseTime=True", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), port)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})

		if err != nil {
			panic("failed to connect database")
		}
		DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.Config("MYSQL_DBNAME")))
		DB.Exec(fmt.Sprintf("USE %s", config.Config("MYSQL_DBNAME")))

		fmt.Printf("Success Create Database: '%s' \n", config.Config("MYSQL_DBNAME"))
		//Migrate the DB
		fmt.Println("Migrating Table...")
		errs := DB.AutoMigrate(&models.Todo{}, &models.Activity{})
		if errs != nil {
			fmt.Sprintln(errs)
		}
		fmt.Println("Table Migrated!")
		sqlDB, err := DB.DB()
		sqlDB.Close()

		if err != nil {
			panic("failed to connect database")
		}
	}
	fmt.Println("Im a Child Process! Connect DB on Child Process")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&parseTime=True", config.Config("MYSQL_USER"), config.Config("MYSQL_PASSWORD"), config.Config("MYSQL_HOST"), port, config.Config("MYSQL_DBNAME"))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		panic("failed to connect database")
	}
}
