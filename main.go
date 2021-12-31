package main

import (
	"todo-apis-go/database"
	"todo-apis-go/router"
	"todo-apis-go/utils"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	utils.NewPool()
	database.ConnectDB()
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	router.Init(app)

	app.Listen(":3030")
}
