package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/restApiOnGo/pkg/domain"
	"github.com/sshaparenko/restApiOnGo/pkg/services"
)

func Signup(c *fiber.Ctx) error {
	//create a var to store the request
	var userInput *domain.UserRequest = new(domain.UserRequest)

	// parse the request into "userInput" variable
	if err := c.BodyParser(userInput); err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}
	// validate the request
	errors := userInput.ValidateStruct()
	// if validation is failed, return the validation errors
	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.Response[[]*domain.ErrorResponse]{
			Success: false,
			Message: "validation failed",
			Data:    errors,
		})
	}
	// perform the signup
	// if signup is successful, the JWT token is returned
	token, err := services.Signup(*userInput)
	// if signup is failed, return an error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}
	// return the JWT token
	return c.JSON(domain.Response[string]{
		Success: true,
		Message: "token data",
		Data:    token,
	})
}

func Login(c *fiber.Ctx) error {
	// create a variable to store the request
	var userInput *domain.UserRequest = new(domain.UserRequest)
	// parse the request into "userInput" variable
	if err := c.BodyParser(userInput); err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}
	// validate the request
	errors := userInput.ValidateStruct()
	// if validation is failed, return the validation errors
	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.Response[[]*domain.ErrorResponse]{
			Success: false,
			Message: "validation failed",
			Data:    errors,
		})
	}
	// perform the login
	// if login is successful, the JWT token is returned
	token, err := services.Login(*userInput)
	// if login is failed, return an error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.Response[any]{
			Success: false,
			Message: err.Error(),
		})
	}
	// return the JWT token
	return c.JSON(domain.Response[string]{
		Success: true,
		Message: "token data",
		Data:    token,
	})
}
