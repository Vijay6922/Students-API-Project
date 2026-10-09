package student

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Vijay6922/StudentsAPIGoProject/internal/types"
	"github.com/Vijay6922/StudentsAPIGoProject/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func NewStudent() http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			if errors.Is(err, io.EOF) {
				err = errors.New("request body is required")
			}
			if writeErr := response.SendResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()}); writeErr != nil {
				log.Printf("failed to write bad request response: %v", writeErr)
			}
			return
		}

		//request validation
		if err := validate.Struct(student); err != nil {
			if writeErr := response.ValidationErrorResponse(w, http.StatusBadRequest, err); writeErr != nil {
				log.Printf("failed to write validation response: %v", writeErr)
			}
			return
		}

		if err := response.SendResponse(w, http.StatusCreated, map[string]string{"success": "Student created"}); err != nil {
			log.Printf("failed to write student response: %v", err)
		}
	}
}
