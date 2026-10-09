package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	StatusCode int    `json:"-"`
	Message    any    `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
}

func (response Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	return json.NewEncoder(w).Encode(response)
}

func SendResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, err error) error {
	return Response{
		StatusCode: statusCode,
		Error:      err.Error(),
	}.Send(w)
}

func SendSuccessResponse(w http.ResponseWriter, statusCode int, data any) error {
	return Response{
		StatusCode: statusCode,
		Message:    data,
	}.Send(w)
}

func ValidationErrorResponse(w http.ResponseWriter, statusCode int, err error) error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return SendErrorResponse(w, statusCode, err)
	}

	var messages []string
	for _, validationErr := range validationErrors {
		switch validationErr.Tag() {
		case "required":
			messages = append(messages, validationErr.Field()+" is required")
		case "email":
			messages = append(messages, validationErr.Field()+" must be a valid email address")
		case "min":
			messages = append(messages, validationErr.Field()+" must be greater than or equal to "+validationErr.Param())
		default:
			messages = append(messages, validationErr.Field()+" is invalid")
		}
	}
	return Response{
		StatusCode: statusCode,
		Message:    messages,
	}.Send(w)
}
