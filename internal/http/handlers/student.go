package student

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Vijay6922/StudentsAPIGoProject/internal/storage"
	"github.com/Vijay6922/StudentsAPIGoProject/internal/types"
	"github.com/Vijay6922/StudentsAPIGoProject/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func NewStudent(store storage.Storage) http.HandlerFunc {
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

		lastID, err := store.CreateStudent(student)
		if err != nil {
			log.Printf("failed to create student: %v", err)
			if writeErr := response.SendResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create student"}); writeErr != nil {
				log.Printf("failed to write internal server error response: %v", writeErr)
			}
			return
		}

		if err := response.SendResponse(w, http.StatusCreated, map[string]any{"success": "Student created", "id": lastID}); err != nil {
			log.Printf("failed to write student response: %v", err)
		}
	}
}

func GetStudent(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || id <= 0 {
			if writeErr := response.SendErrorResponse(w, http.StatusBadRequest, errors.New("student ID must be a positive integer")); writeErr != nil {
				log.Printf("failed to write invalid student ID response: %v", writeErr)
			}
			return
		}

		student, err := store.GetStudent(id)
		if errors.Is(err, storage.ErrStudentNotFound) {
			if writeErr := response.SendErrorResponse(w, http.StatusNotFound, err); writeErr != nil {
				log.Printf("failed to write student not found response: %v", writeErr)
			}
			return
		}
		if err != nil {
			log.Printf("failed to get student %d: %v", id, err)
			if writeErr := response.SendErrorResponse(w, http.StatusInternalServerError, errors.New("failed to get student")); writeErr != nil {
				log.Printf("failed to write internal server error response: %v", writeErr)
			}
			return
		}

		if err := response.SendResponse(w, http.StatusOK, student); err != nil {
			log.Printf("failed to write student response: %v", err)
		}
	}
}
