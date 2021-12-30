package main

import (
	"encoding/json"
	"todo-apis-go/database"
	"todo-apis-go/router"

	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to Todo API!",
	})
}

func main() {
	database.ConnectDB()
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		AppName:     "Todo-API",
		Prefork:     true,
	})

	app.Get("/", welcome)

	router.Init(app)

	app.Listen(":3030")
}
