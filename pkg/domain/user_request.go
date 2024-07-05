package domain

import "github.com/go-playground/validator/v10"

type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (userInput UserRequest) ValidateStruct() []*ErrorResponse {
	var errors []*ErrorResponse

	validate := validator.New()
	err := validate.Struct(userInput)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.ErrorMessage = getErrorMessage(err)
			element.Field = err.Field()
			errors = append(errors, &element)
		}
	}
	return errors
}
