package models

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	Field        string `json:"field"`
}
