package router

import (
	"fmt"
	"todo-apis-go/cache"
	"todo-apis-go/services/activity"
	"todo-apis-go/services/todo"

	"github.com/gofiber/fiber/v2"
)

func cacheMiddleware(c *fiber.Ctx) error {
	id := c.Params("id")

	activityId := c.Query("activity_group_id")

	if id == "0" {
		if activityId != "" {
			id = fmt.Sprintf("agi_%s", activityId)
		} else {
			id = "all"
		}
	}
	result, err := cache.Cache.Get(id)
	if err == nil {
		return c.JSON(fiber.Map{"data": result, "status": "Success", "message": "Success"})
	} else {
		if activityId != "" {
			return c.JSON(fiber.Map{
				"status":  "Success",
				"message": "Success",
				"data":    make([]string, 0),
			})

		}
		return c.Next()
	}
}

func Init(entrypoint *fiber.App) {
	// entrypoint.Use(logger.New())
	entrypoint.Get("/todo-items", cacheMiddleware, todo.GetAll)
	entrypoint.Get("/todo-items/:id", cacheMiddleware, todo.GetById)
	entrypoint.Post("/todo-items", todo.Store)
	entrypoint.Delete("/todo-items/:id", todo.Destroy)
	entrypoint.Patch("/todo-items/:id", todo.Update)
	entrypoint.Get("/activity-groups", activity.GetAll)
	entrypoint.Get("/activity-groups/:id", activity.GetById)
	entrypoint.Post("/activity-groups", activity.Store)
	entrypoint.Delete("/activity-groups/:id", activity.Destroy)
	entrypoint.Patch("/activity-groups/:id", activity.Update)
}
