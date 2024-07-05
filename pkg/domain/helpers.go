package domain

import "github.com/go-playground/validator/v10"

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "gt":
		return "the velue of " + err.Field() + " mus be greater then " + err.Param()
	case "gte":
		return "the velue of " + err.Field() + " mus be greater or equal to " + err.Param()
	case "email":
		return "the email is invalid"
	case "min":
		return "the minimum length of " + err.Field() + " equals to " + err.Param()
	default:
		return "validation error in " + err.Field()
	}
}
