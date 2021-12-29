package router

import (
	"todo-apis-go/services/activity"
	"todo-apis-go/services/todo"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func middleware(app *fiber.App) {
	app.Use(etag.New())
	// app.Use(compress.New(compress.Config{
	// 	Level: compress.LevelBestSpeed, // 1
	// }))
}

func Init(app *fiber.App) {
	// middleware(app)

	endpoint := app.Group("/", logger.New())
	endpoint.Get("/todo-items", todo.GetAll)
	endpoint.Get("/todo-items/:id", todo.GetById)
	endpoint.Post("/todo-items", todo.Store)
	endpoint.Delete("/todo-items/:id", todo.Destroy)
	endpoint.Patch("/todo-items/:id", todo.Update)
	endpoint.Put("/todo-items/:id", todo.Update)

	endpoint.Get("/activity-groups", activity.GetAll)
	endpoint.Get("/activity-groups/:id", activity.GetById)
	endpoint.Post("/activity-groups", activity.Store)
	endpoint.Delete("/activity-groups/:id", activity.Destroy)
	endpoint.Patch("/activity-groups/:id", activity.Update)
	endpoint.Put("/activity-groups/:id", activity.Update)
}
