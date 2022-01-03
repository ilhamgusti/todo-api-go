package activity

import (
	"fmt"
	"strconv"
	"todo-apis-go/cache"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/eko/gocache/store"
	"github.com/gofiber/fiber/v2"
)

var cacheActivityId uint

type EmptyMap struct{}

func GetAll(c *fiber.Ctx) error {
	var result interface{}
	var err error

	go func() {
		var activities []models.Activity
		err = database.DB.Find(&activities).Error
		if err != nil {
			cache.Cache.Delete("activities")
		}
		cache.Cache.Set("activities", activities, &store.Options{Cost: 1})
	}()

	result, err = cache.Cache.Get("activities")
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    make([]string, 0),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})
}

func GetById(c *fiber.Ctx) error {

	var err error
	var result interface{}
	id := c.Params("id")
	// id, err = c.ParamsInt("id")
	// if err != nil {
	// 	return c.Status(404).JSON(fiber.Map{
	// 		"status":  "Not Found",
	// 		"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
	// 		"data":    EmptyMap{},
	// 	})
	// }
	//check inside cache
	result, err = cache.Cache.Get("a" + id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}

	if result != nil {
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    result,
		})
	}

	go func() {
		var activity models.Activity
		err = database.DB.First(&activity, id).Error
		if err != nil {
			cache.Cache.Delete("a" + id)
		}
		cache.Cache.Set("a"+id, activity, &store.Options{Cost: 1})

	}()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
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
	var err error
	cacheActivityId = cacheActivityId + 1

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

	activity.ID = cacheActivityId

	go func() {
		database.DB.Create(&activity)
		err = cache.Cache.Set("a"+strconv.Itoa(int(activity.ID)), &activity, &store.Options{Cost: 1})
	}()
	err = cache.Cache.Set("a"+strconv.Itoa(int(activity.ID)), &activity, &store.Options{Cost: 1})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Activity with ID %d Not Save To Cache`, &activity.ID),
			"data":    EmptyMap{},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    activity,
	})
}

func Destroy(c *fiber.Ctx) error {
	var err error

	id := c.Params("id")

	//check inside cache
	_, err = cache.Cache.Get("a" + id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %s Not Found`, id),
			"data":    EmptyMap{},
		})
	}
	err = cache.Cache.Delete("a" + id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Activity with ID %s Not Save To Cache`, id),
			"data":    EmptyMap{},
		})
	}
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
	var id string
	var err error
	var result interface{}

	id = c.Params("id")

	result, err = cache.Cache.Get("a" + id) //, new(models.Activity)
	if err != nil {
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
	}()

	err = cache.Cache.Set("a"+id, result, &store.Options{Cost: 1})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Activity with ID %s Not Save To Cache`, id),
			"data":    EmptyMap{},
		})
	}
	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
