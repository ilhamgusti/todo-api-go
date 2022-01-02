package activity

import (
	"fmt"
	"todo-apis-go/cache"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/eko/gocache/store"
	"github.com/gofiber/fiber/v2"
)

type KeyOfActivity struct {
	ID int
}

var cacheActivityId int

func GetAll(c *fiber.Ctx) error {
	var result interface{}
	var err error

	go func() {
		var activities []models.Activity
		err = database.DB.Find(&activities).Error
		if err != nil {
			cache.Cache.Delete("activities")
		}
		cache.Cache.Set("activities", activities, &store.Options{Cost: 4})
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
	var id int
	var err error
	var result interface{}

	id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}
	//check inside cache
	result, err = cache.Cache.Get(KeyOfActivity{ID: id})
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    make(map[string]string),
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
			cache.Cache.Delete(KeyOfActivity{ID: id})
		}
		cache.Cache.Set(KeyOfActivity{ID: id}, activity, &store.Options{Cost: 4})

	}()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    make(map[string]string),
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
		err = cache.Cache.Set(KeyOfActivity{ID: int(activity.ID)}, &activity, &store.Options{Cost: 4})
	}()
	err = cache.Cache.Set(KeyOfActivity{ID: int(activity.ID)}, &activity, &store.Options{Cost: 4})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Activity with ID %d Not Save To Cache`, &activity.ID),
			"data":    make(map[string]string),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    activity,
	})
}

func Destroy(c *fiber.Ctx) error {
	var id int
	var err error

	id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	//check inside cache
	_, err = cache.Cache.Get(KeyOfActivity{ID: id})
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}
	err = cache.Cache.Delete(KeyOfActivity{ID: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Activity with ID %d Not Save To Cache`, id),
			"data":    make(map[string]string),
		})
	}
	go func() {
		database.DB.Delete(&models.Activity{}, id)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    make(map[string]string),
	})
}

func Update(c *fiber.Ctx) error {
	var id int
	var err error
	var result interface{}

	id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	result, err = cache.Cache.Get(KeyOfActivity{ID: id}) //, new(models.Activity)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	if err := c.BodyParser(&result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": fmt.Sprintf(`Activity with ID %d Bad Request`, id),
		})
	}

	go func() {
		database.DB.Save(&result)
	}()

	err = cache.Cache.Set(KeyOfActivity{ID: int(id)}, result, &store.Options{Cost: 4})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Activity with ID %d Not Save To Cache`, id),
			"data":    make(map[string]string),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
