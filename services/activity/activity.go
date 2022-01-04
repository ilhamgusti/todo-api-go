package activity

import (
	"fmt"
	"strconv"
	"sync"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/gofiber/fiber/v2"
)

var cacheActivityId uint

type EmptyMap struct{}

var C sync.Map

func GetAll(c *fiber.Ctx) error {

	go func() {
		var activities []models.Activity
		err := database.DB.Find(&activities).Limit(1).Error
		if err != nil {
			C.Delete("activities")
		}
		C.Store("activities", activities)
	}()

	result, ok := C.Load("activities")
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
	result, ok := C.Load("a" + id)
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
		err := database.DB.First(&activity, id).Error
		if err != nil {
			C.Delete("a" + id)
		}
		C.Store("a"+id, activity)

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "Bad Request",
		})
	}
	if activity.Title == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "title cannot be null",
		})
	}
	if activity.Email == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "email cannot be null",
		})
	}

	activity.ID = cacheActivityId + 1
	cacheActivityId = activity.ID
	go func() {
		database.DB.Select("Title", "Email").Create(&activity)
		C.Store("a"+strconv.Itoa(int(activity.ID)), &activity)
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
	result, ok := C.Load("a" + id)
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
	C.Delete("a" + id)

	go func() {
		database.DB.Unscoped().Delete(&models.Activity{}, id)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    EmptyMap{},
	})
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")

	result, ok := C.Load("a" + id) //, new(models.Activity)
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	if err := c.BodyParser(&result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": fmt.Sprintf(`Activity with ID %s Bad Request`, id),
		})
	}

	go func() {
		database.DB.Save(&result)
		C.Store("a"+id, result)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
