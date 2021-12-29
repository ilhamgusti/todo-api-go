package todo

import (
	"errors"
	"fmt"
	"strconv"
	"todo-apis-go/database"
	"todo-apis-go/models"
	"todo-apis-go/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EmptyMap struct{}

func GetAll(c *fiber.Ctx) error {
	db := database.DB
	var todos []models.Todo
	activityId := c.Query("activity_group_id")

	if activityId == "" {
		db.Table("todos").Find(&todos)
		utils.Cache.Set("all", todos)
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    &todos,
		})
	}
	db.Raw("SELECT * FROM todos WHERE activity_group_id = ?", activityId).Scan(&todos)

	utils.Cache.Set(fmt.Sprintf("agi_%s", activityId), todos)

	if todos != nil {
		todos = []models.Todo{}
	}
	fmt.Println(todos)

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    &todos,
	})

}

func GetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	var todo models.Todo

	errdb := database.DB.First(&todo, id).Error

	if errors.Is(errdb, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	utils.Cache.Set(strconv.Itoa(id), todo)

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    &todo,
	})

}

func Store(c *fiber.Ctx) error {
	todo := new(models.Todo)

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "Bad Request",
		})
	}
	if todo.Title == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "title cannot be null",
		})
	}
	if todo.ActivityGroupId == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "activity_group_id cannot be null",
		})
	}

	todo.IsActive = true
	todo.Priority = "very-high"

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "Bad Request",
		})
	}
	database.DB.Create(&todo)

	return c.Status(201).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    &todo,
	})
}

func Destroy(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	result := database.DB.First(&models.Todo{}, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	result.Delete(&models.Todo{})

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    EmptyMap{},
	})
}

func Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	todo := new(models.Todo)

	errdb := database.DB.First(&todo, id).Error

	if errors.Is(errdb, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": fmt.Sprintf(`Todo with ID %d Bad Request`, id),
		})
	}

	database.DB.Save(&todo)
	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    &todo,
	})

}
