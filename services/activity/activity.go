package activity

import (
	"fmt"
	"strconv"
	"sync"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/gofiber/fiber/v2"
)

var CacheActivityId int

type EmptyMap struct{}

var Cache sync.Map

func GetAll(c *fiber.Ctx) error {
	go func() {
		var activities []models.Activity
		if err := database.DB.Find(&activities).Limit(1).Error; err != nil {
			Cache.Delete("activities")
		}
		Cache.Store("activities", activities)
	}()

	result, ok := Cache.Load("activities")
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

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})
}

func GetById(c *fiber.Ctx) error {
	id := c.Params("id")
	//check inside Cache
	result, ok := Cache.Load("a" + id)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	if result == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	go func() {
		var activity models.Activity
		if err := database.DB.First(&activity, id).Error; err != nil {
			Cache.Delete("a" + id)
		}
		Cache.Store("a"+id, activity)

	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})
}

func Store(c *fiber.Ctx) error {
	activity := new(models.Activity)

	if err := c.BodyParser(&activity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "Bad Request", "message": "body parser error"})
	}
	if activity.Title == nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "Bad Request", "message": "title cannot be null"})
	}
	if activity.Email == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "Bad Request", "message": "email cannot be null"})
	}
	CacheActivityId = CacheActivityId + 1
	activity.ID = CacheActivityId
	go func() {
		database.DB.Select("Title", "Email").Create(&activity)
		Cache.Store("a"+strconv.Itoa(activity.ID), activity)
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    activity,
	})
}

func Destroy(c *fiber.Ctx) error {
	id := c.Params("id")

	//check inside Cache
	result, ok := Cache.Load("a" + id)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	if result == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	Cache.Delete("a" + id)

	go func() {
		database.DB.Delete(&models.Activity{}, id)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    EmptyMap{},
	})
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")

	result, ok := Cache.Load("a" + id)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	c.BodyParser(&result)

	go func() {
		database.DB.Save(&result)
		Cache.Store("a"+id, result)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
