package activity

import (
	"errors"
	"fmt"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EmptyMap struct {
}

func GetAll(c *fiber.Ctx) error {
	var activities []models.Activity

	database.DB.Table("activities").Find(&activities)
	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    &activities,
	})
}

func GetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	var activity models.Activity

	errdb := database.DB.First(&activity, id).Error

	if errors.Is(errdb, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    &activity,
	})

}
func Store(c *fiber.Ctx) error {

	activity := new(models.Activity)

	// errorValidate := validators.ValidateActivityStruct(*activity)
	// fmt.Println(errorValidate)
	// if errorValidate != nil {
	// 	return c.Status(400).JSON(errorValidate)
	// }
	if err := c.BodyParser(&activity); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "Bad Request",
		})
	}
	if activity.Title == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "title cannot be null",
		})
	}
	if activity.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "email cannot be null",
		})
	}
	database.DB.Create(&activity)

	return c.Status(201).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    activity,
	})
}

func Destroy(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	var activity models.Activity
	result := database.DB.First(&activity, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	result.Delete(&activity)

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
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	activity := new(models.Activity)

	errdb := database.DB.First(&activity, id).Error

	if errors.Is(errdb, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	if err := c.BodyParser(&activity); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": fmt.Sprintf(`Activity with ID %d Bad Request`, id),
		})
	}

	// errorValidate := validators.ValidateActivityStruct(*activity)
	// fmt.Println(errorValidate)
	// if errorValidate != nil {
	// 	return c.Status(400).JSON(errorValidate)
	// }

	database.DB.Save(&activity)
	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    &activity,
	})

}
