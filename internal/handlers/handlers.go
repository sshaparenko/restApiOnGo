package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/restApiOnGo/internal/models"
	"github.com/sshaparenko/restApiOnGo/internal/services"
	"github.com/sshaparenko/restApiOnGo/internal/utils"
)

func GetAllItems(c * fiber.Ctx) error {
	var items []models.Item = services.GetAllItems()

	return c.JSON(models.Response[[]models.Item]{
		Success: true,
		Message: "All items data",
		Data: items,
	})
}

func GetItemByID(c * fiber.Ctx) error {
	var itemID string = c.Params("id")

	item, err := services.GetItemById(itemID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(models.Response[models.Item]{
		Success: true,
		Message: "item found",
		Data: item,
	})
}

func CreateItem(c * fiber.Ctx) error {
	isValid, err := utils.CheckToken(c)

	if !isValid {
		return c.Status(http.StatusUnauthorized).JSON(models.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	//create variable to store the request
	var itemInput *models.ItemRequest = new(models.ItemRequest)

	//parse the request into "itemInput" variable
	if err := c.BodyParser(itemInput); err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}
	//validate the request
	errors := itemInput.ValidateStruct()

	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response[[]*models.ErrorResponse]{
			Success: false,
			Message: "validation failed",
			Data: errors,
		})
	}
	//create a new item from validated request
	var createItem models.Item = services.CreateItem(*itemInput)

	return c.Status(http.StatusCreated).JSON(models.Response[models.Item]{
		Success: true,
		Message: "item created",
		Data: createItem,
	})
}

func UpdateItem(c * fiber.Ctx) error {
	isValid, err := utils.CheckToken(c)

	if !isValid {
		return c.Status(http.StatusUnauthorized).JSON(models.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	var itemInput *models.ItemRequest = new(models.ItemRequest)

	if err := c.BodyParser(itemInput); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	errors := itemInput.ValidateStruct()

	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(models.Response[[]*models.ErrorResponse]{
			Success: false,
			Message: "valdation failed",
			Data: errors,
		})
	}

	var itemID string = c.Params("id")

	updatedItem, err := services.UpdateItem(*itemInput, itemID)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(models.Response[models.Item]{
		Success: true,
		Message: "item updated",
		Data: updatedItem,
	})
}

func DeleteItem(c *fiber.Ctx) error {
	isValid, err := utils.CheckToken(c)

	if !isValid {
		return c.Status(http.StatusUnauthorized).JSON(models.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	var itemId string = c.Params("id")

	var result = services.DeleteItem(itemId)

	if result {
		return c.JSON(models.Response[any]{
			Success: true,
			Message: "item deleted",
		})
	}

	return c.Status(http.StatusNotFound).JSON(models.Response[any]{
		Success: false,
		Message: "item failed to delete",
	})
}