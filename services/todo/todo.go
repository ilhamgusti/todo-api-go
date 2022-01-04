package todo

import (
	"fmt"
	"strconv"
	"sync"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/gofiber/fiber/v2"
)

var cacheID uint

var C sync.Map

type EmptyMap struct{}

func GetAll(c *fiber.Ctx) error {

	id := c.Query("activity_group_id")

	if id == "" {
		result, ok := C.Load("all")
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
			C.Store("all", &todos)
		}()

		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    result,
		})
	}

	result, ok := C.Load("agi_" + id)
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
		C.Store("agi_"+id, &todos)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}

func GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, ok := C.Load("t" + id)
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
			C.Store("t"+id, &todo)
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

	todo.IsActive = true
	todo.Priority = "very-high"
	todo.ID = cacheID + 1
	cacheID = todo.ID

	go func() {
		database.DB.Create(&todo)
		C.Store("t"+strconv.Itoa(int(todo.ID)), &todo)
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
	_, ok := C.Load("t" + id) //, new(models.Todo)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	go func() {
		database.DB.Unscoped().Delete(&models.Todo{}, id)
	}()
	C.Delete("t" + id)

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    EmptyMap{},
	})
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")

	result, ok := C.Load("t" + id) // , new(models.Todo)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	if err := c.BodyParser(&result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": fmt.Sprintf(`Todo with ID %s Bad Request`, id),
		})
	}

	go func() {
		database.DB.Save(&result)
		C.Store("t"+id, result)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
