package activity

import (
	"fmt"
	"strconv"
	"sync"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/gofiber/fiber/v2"
)

var cacheActivityId int

type EmptyMap struct{}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var cache sync.Map

func GetAll(c *fiber.Ctx) error {
	go func() {
		var activities []models.Activity
		if err := database.DB.Find(&activities).Limit(1).Error; err != nil {
			cache.Delete("activities")
		}
		cache.Store("activities", activities)
	}()

	result, ok := cache.Load("activities")
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
	//check inside cache
	result, ok := cache.Load(id)
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
			cache.Delete(id)
		}
		cache.Store(id, activity)

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
		return c.Status(fiber.StatusBadRequest).JSON(Response{Status: "Bad Request", Message: "body parser error"})
	}
	if activity.Title == nil {

		return c.Status(fiber.StatusBadRequest).JSON(Response{Status: "Bad Request", Message: "title cannot be null"})
	}
	if activity.Email == nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Status: "Bad Request", Message: "email cannot be null"})
	}
	cacheActivityId = cacheActivityId + 1
	activity.ID = strconv.Itoa(cacheActivityId)
	go func() {
		database.DB.Select("Title", "Email").Create(&activity)
		// cache.Store(activity.ID, activity)
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    activity,
	})
}

func Destroy(c *fiber.Ctx) error {
	id := c.Params("id")

	//check inside cache
	result, ok := cache.Load(id)
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
	cache.Delete(id)

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

	result, ok := cache.Load(id)
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
		cache.Store(id, result)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
