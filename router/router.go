package router

import (
	"fmt"
	"strconv"
	"todo-apis-go/services/activity"
	"todo-apis-go/services/todo"
	"todo-apis-go/utils"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/etag"
)

func cacheMiddleware(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	stringId := strconv.Itoa(id)

	activityId := c.Query("activity_group_id")

	if stringId == "0" {
		if activityId != "" {
			stringId = fmt.Sprintf("agi_%s", activityId)
		} else {
			stringId = "all"
		}
	}
	val, err := utils.Cache.Get(stringId)
	if err != ttlcache.ErrNotFound {
		return c.JSON(fiber.Map{"data": val, "status": "Success", "message": "Success"})
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

func middleware(app *fiber.App) {
	app.Use(etag.New())
	app.Use(cache.New())
}

func Init(app *fiber.App) {
	// middleware(app)

	endpoint := app.Group("/")
	endpoint.Get("/todo-items", cacheMiddleware, todo.GetAll)
	endpoint.Get("/todo-items/:id", cacheMiddleware, todo.GetById)
	endpoint.Post("/todo-items", todo.Store)
	endpoint.Delete("/todo-items/:id", todo.Destroy)
	endpoint.Patch("/todo-items/:id", todo.Update)

	endpoint.Get("/activity-groups", activity.GetAll)
	endpoint.Get("/activity-groups/:id", activity.GetById)
	endpoint.Post("/activity-groups", activity.Store)
	endpoint.Delete("/activity-groups/:id", activity.Destroy)
	endpoint.Patch("/activity-groups/:id", activity.Update)
}
