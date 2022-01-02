package todo

import (
	"fmt"
	"todo-apis-go/cache"
	"todo-apis-go/database"
	"todo-apis-go/models"

	"github.com/eko/gocache/store"
	"github.com/gofiber/fiber/v2"
)

type KeyOfTodo struct {
	ID int
}

var cacheID int

func GetAll(c *fiber.Ctx) error {
	var id string
	id = c.Query("activity_group_id")
	var todos []models.Todo

	if id == "" {
		database.DB.Table("todos").Find(&todos)
		cache.Cache.Set("all", todos, &store.Options{Cost: 4})
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    todos,
		})
	}
	database.DB.Where("activity_group_id = ?", id).Find(&todos)

	if todos != nil {
		cache.Cache.Set(fmt.Sprintf("agi_%s", id), todos, &store.Options{Cost: 4})
		return c.JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    make([]string, 0),
		})

	}
	cache.Cache.Set(fmt.Sprintf("agi_%s", id), todos, &store.Options{Cost: 4})

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    todos,
	})

}

func GetById(c *fiber.Ctx) error {
	var err error
	var id int

	id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Activity with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	//check inside cache
	result, err := cache.Cache.Get(KeyOfTodo{ID: id}) //, new(models.Todo)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	go func() {
		var todo models.Todo
		err = database.DB.First(&todo, id).Error
		err = cache.Cache.Set(KeyOfTodo{ID: id}, todo, &store.Options{Cost: 4})
	}()

	//save to cache for future check
	err = cache.Cache.Set(KeyOfTodo{ID: id}, result, &store.Options{Cost: 4})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Todo with ID %d Not Save To Cache`, id),
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

	cacheID = cacheID + 1

	todo := new(models.Todo)

	if err = c.BodyParser(&todo); err != nil {
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
	todo.ID = cacheID

	go func() {
		database.DB.Create(&todo)
		err = cache.Cache.Set(KeyOfTodo{ID: int(todo.ID)}, &todo, &store.Options{Cost: 4})
	}()

	err = cache.Cache.Set(KeyOfTodo{ID: int(todo.ID)}, &todo, &store.Options{Cost: 4})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Todo with ID %d Not Save To Cache`, todo.ID),
			"data":    make(map[string]string),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    todo,
	})
}

func Destroy(c *fiber.Ctx) error {
	var err error
	var id int

	id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	//check inside cache
	_, err = cache.Cache.Get(KeyOfTodo{ID: id}) //, new(models.Todo)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	err = cache.Cache.Delete(KeyOfTodo{ID: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": fmt.Sprintf(`Todo with ID %d Not Save To Cache`, id),
			"data":    make(map[string]string),
		})
	}

	go func() {
		database.DB.Unscoped().Delete(&models.Todo{}, id)
	}()

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    make(map[string]string),
	})
}

func Update(c *fiber.Ctx) error {
	var err error
	var id int

	id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	result, err := cache.Cache.Get(KeyOfTodo{ID: id}) // , new(models.Todo)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "Not Found",
			"message": fmt.Sprintf(`Todo with ID %d Not Found`, id),
			"data":    make(map[string]string),
		})
	}

	if err := c.BodyParser(&result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": fmt.Sprintf(`Todo with ID %d Bad Request`, id),
		})
	}

	go func() {
		database.DB.Save(&result)
	}()

	err = cache.Cache.Set(KeyOfTodo{ID: int(id)}, &result, &store.Options{Cost: 4})

	return c.JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    result,
	})

}
