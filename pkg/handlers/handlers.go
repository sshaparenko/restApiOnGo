package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/restApiOnGo/pkg/domain"
	"github.com/sshaparenko/restApiOnGo/pkg/services"
	"github.com/sshaparenko/restApiOnGo/pkg/utils"
)

func GetAllItems(c *fiber.Ctx) error {
	items, err := services.GetAllItems()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.Response[[]domain.Item]{
		Success: true,
		Message: "All items data",
		Data:    items,
	})
}

func GetItemByID(c *fiber.Ctx) error {

	var itemID string = c.Params("id")

	item, err := services.GetItemByID(itemID)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.Response[domain.Item]{
		Success: true,
		Message: "item found",
		Data:    item,
	})
}

func CreateItem(c *fiber.Ctx) error {
	isValid, err := utils.CheckToken(c)

	if !isValid {
		return c.Status(http.StatusUnauthorized).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	//create variable to store the request
	var itemInput *domain.ItemRequest = new(domain.ItemRequest)

	//parse the request into "itemInput" variable
	if err := c.BodyParser(itemInput); err != nil {
		return c.Status(http.StatusNotFound).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}
	//validate the request
	errors := itemInput.ValidateStruct()

	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.Response[[]*domain.ErrorResponse]{
			Success: false,
			Message: "validation failed",
			Data:    errors,
		})
	}

	item, err := services.CreateItem(*itemInput)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.Response[domain.Item]{
			Success: false,
			Message: err.Error(),
			Data:    domain.Item{},
		})
	}

	return c.Status(http.StatusCreated).JSON(domain.Response[domain.Item]{
		Success: true,
		Message: "item created",
		Data:    item,
	})
}

func UpdateItem(c *fiber.Ctx) error {
	isValid, err := utils.CheckToken(c)

	if !isValid {
		return c.Status(http.StatusUnauthorized).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	var itemInput *domain.ItemRequest = new(domain.ItemRequest)

	if err := c.BodyParser(itemInput); err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	errors := itemInput.ValidateStruct()

	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.Response[[]*domain.ErrorResponse]{
			Success: false,
			Message: "valdation failed",
			Data:    errors,
		})
	}

	var itemID string = c.Params("id")

	updatedItem, err := services.UpdateItem(*itemInput, itemID)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.Response[domain.Item]{
		Success: true,
		Message: "item updated",
		Data:    updatedItem,
	})
}

func DeleteItem(c *fiber.Ctx) error {
	isValid, err := utils.CheckToken(c)

	if !isValid {
		return c.Status(http.StatusUnauthorized).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	var itemID string = c.Params("id")

	result, err := services.DeleteItem(itemID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.Response[domain.Item]{
			Success: false,
			Message: err.Error(),
			Data:    domain.Item{},
		})
	}

	if result {
		return c.JSON(domain.Response[any]{
			Success: true,
			Message: "item deleted",
		})
	}

	return c.Status(http.StatusNotFound).JSON(domain.Response[any]{
		Success: false,
		Message: "item failed to delete",
	})
}
