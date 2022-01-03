package todo

import (
	"fmt"
	"strconv"
	"todo-apis-go/cache"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/eko/gocache/store"
	"github.com/gofiber/fiber/v2"
)

var cacheID uint

type EmptyMap struct{}

func GetAll(c *fiber.Ctx) error {

	id := c.Query("activity_group_id")
	var todos []models.Todo

	if id == "" {
		database.DB.Table("todos").Find(&todos)
		cache.Cache.Set("all", todos, &store.Options{Cost: 2})
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    todos,
		})
	}
	database.DB.Where("activity_group_id = ?", id).Find(&todos)

	if todos != nil {
		cache.Cache.Set(fmt.Sprintf("agi_%s", id), todos, &store.Options{Cost: 2})
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    make([]string, 0),
		})
	}
	cache.Cache.Set(fmt.Sprintf("agi_%s", id), todos, &store.Options{Cost: 2})

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    todos,
	})

}

func GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := cache.Cache.Get("t" + id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	go func() {
		var todo models.Todo
		err = database.DB.First(&todo, id).Error
	}()

	//save to cache for future check
	err = cache.Cache.Set("t"+id, result, &store.Options{Cost: 2})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Todo with ID %s Not Save To Cache`, id),
			"data":    EmptyMap{},
		})
	}
	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}

func Store(c *fiber.Ctx) error {
	todo := new(models.Todo)

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "Bad Request",
		})
	}
	if todo.Title == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "title cannot be null",
		})
	}
	if todo.ActivityGroupId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "activity_group_id cannot be null",
		})
	}
	cacheID = cacheID + 1

	todo.IsActive = true
	todo.Priority = "very-high"
	todo.ID = cacheID

	go func() {
		database.DB.Create(&todo)
		cache.Cache.Set("t"+strconv.Itoa(int(todo.ID)), &todo, &store.Options{Cost: 2})
	}()

	err := cache.Cache.Set("t"+strconv.Itoa(int(todo.ID)), &todo, &store.Options{Cost: 2})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Todo with ID %d Not Save To Cache`, todo.ID),
			"data":    EmptyMap{},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    todo,
	})
}

func Destroy(c *fiber.Ctx) error {
	id := c.Params("id")

	//check inside cache
	_, err := cache.Cache.Get("t" + id) //, new(models.Todo)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	err = cache.Cache.Delete("t" + id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Todo with ID %s Not Save To Cache`, id),
			"data":    EmptyMap{},
		})
	}

	go func() {
		database.DB.Unscoped().Delete(&models.Todo{}, id)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    EmptyMap{},
	})
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := cache.Cache.Get("t" + id) // , new(models.Todo)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	if err = c.BodyParser(&result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": fmt.Sprintf(`Todo with ID %s Bad Request`, id),
		})
	}

	go func() {
		database.DB.Save(&result)
	}()

	cache.Cache.Set("t"+id, result, &store.Options{Cost: 2})

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
