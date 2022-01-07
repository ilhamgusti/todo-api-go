package router

import (
	"todo-apis-go/services/activity"
	"todo-apis-go/services/todo"

	"github.com/gofiber/fiber/v2"
)

func Init(entrypoint *fiber.App) {
	entrypoint.Get("/todo-items", todo.GetAll)
	entrypoint.Get("/todo-items/:id", todo.GetById)
	entrypoint.Post("/todo-items", todo.Store)
	entrypoint.Delete("/todo-items/:id", todo.Destroy)
	entrypoint.Patch("/todo-items/:id", todo.Update)
	entrypoint.Get("/activity-groups", activity.GetAll)
	entrypoint.Get("/activity-groups/:id", activity.GetById)
	entrypoint.Post("/activity-groups", activity.Store)
	entrypoint.Delete("/activity-groups/:id", activity.Destroy)
	entrypoint.Patch("/activity-groups/:id", activity.Update)
}
