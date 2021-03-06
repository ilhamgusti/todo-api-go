package todo

import (
	"fmt"
	"strconv"
	"sync"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/gofiber/fiber/v2"
)

var cacheID int

var Cache sync.Map

type EmptyMap struct{}

func GetAll(c *fiber.Ctx) error {

	id := c.Query("activity_group_id", "all")

	if id != "all" {
		result, ok := Cache.Load("agi_" + id)
		if !ok {
			return c.JSON(fiber.Map{
				"status":  "Success",
				"message": "Success",
				"data":    []EmptyMap{},
			})
		}
		if result == nil {
			return c.JSON(fiber.Map{
				"status":  "Success",
				"message": "Success",
				"data":    []EmptyMap{},
			})
		}
		go func() {
			var todos []models.Todo
			database.DB.Where("activity_group_id = ?", id).Find(&todos).Limit(1)
			Cache.Store("agi_"+id, todos)
		}()

		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    result,
		})
	}

	result, ok := Cache.Load("t" + id)
	if !ok {
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    []EmptyMap{},
		})
	}
	if result == nil {
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    []EmptyMap{},
		})
	}

	go func() {
		var todos []models.Todo
		database.DB.Table("todos").Find(&todos).Limit(1)
		Cache.Store("all", todos)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}

func GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, ok := Cache.Load("t" + id)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	go func() {
		var todo models.Todo
		err := database.DB.First(&todo, id).Error
		if err != nil {
			Cache.Store("t"+id, todo)
		}
	}()

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
	todo.ID = cacheID
	todo.IsActive = true
	todo.Priority = "very-high"

	go func() {
		database.DB.Create(&todo)
		Cache.Store("t"+strconv.Itoa(todo.ID), todo)
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    todo,
	})
}

func Destroy(c *fiber.Ctx) error {
	id := c.Params("id")

	//check inside cache
	_, ok := Cache.Load("t" + id) //, new(models.Todo)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	Cache.Delete("t" + id)

	go func() {
		database.DB.Delete(&models.Todo{}, id)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    EmptyMap{},
	})
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")

	result, ok := Cache.Load("t" + id)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	c.BodyParser(&result)

	go func() {
		database.DB.Save(&result)
		Cache.Store("t"+id, result)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
