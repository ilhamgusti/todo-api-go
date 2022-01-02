package main

import (
	"todo-apis-go/cache"
	"todo-apis-go/database"
	"todo-apis-go/router"

	"github.com/gofiber/fiber/v2"
	// jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/encoding/json"
	// "encoding/json"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	database.ConnectDB()
	cache.Init()

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
