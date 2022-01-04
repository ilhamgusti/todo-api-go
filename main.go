package main

import (
	"todo-apis-go/database"
	"todo-apis-go/router"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/encoding/json"
)

func main() {
	database.ConnectDB()

	var app *fiber.App = fiber.New(fiber.Config{
		JSONEncoder:                  json.Marshal,
		JSONDecoder:                  json.Unmarshal,
		DisableStartupMessage:        true,
		DisableDefaultDate:           true,
		DisableHeaderNormalizing:     true,
		DisablePreParseMultipartForm: true,
		DisableDefaultContentType:    true,
	})

	router.Init(app)

	app.Listen(":3030")
}
