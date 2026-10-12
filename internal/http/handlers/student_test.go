package student

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vijay6922/StudentsAPIGoProject/internal/storage"
	"github.com/Vijay6922/StudentsAPIGoProject/internal/types"
)

type studentStore struct {
	student types.Student
	err     error
}

func (s studentStore) CreateStudent(types.Student) (int64, error) {
	return 0, nil
}

func (s studentStore) GetStudent(int64) (types.Student, error) {
	return s.student, s.err
}

func TestGetStudent(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		store      studentStore
		wantStatus int
	}{
		{
			name:       "returns student",
			path:       "/api/students/12",
			store:      studentStore{student: types.Student{Id: 12, Name: "Ada Lovelace", Email: "ada@example.com", Age: 36, Phone: "+1 555 0100"}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "rejects invalid ID",
			path:       "/api/students/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "returns not found",
			path:       "/api/students/12",
			store:      studentStore{err: storage.ErrStudentNotFound},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "returns internal error",
			path:       "/api/students/12",
			store:      studentStore{err: errors.New("database unavailable")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("GET /api/students/{id}", GetStudent(tt.store))

			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("GET %s status = %d, want %d; body: %s", tt.path, recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if tt.wantStatus == http.StatusOK {
				if got := recorder.Header().Get("Content-Type"); got != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", got)
				}
				var got types.Student
				if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
					t.Fatalf("decoding student response: %v", err)
				}
				if got != tt.store.student {
					t.Errorf("response student = %+v, want %+v", got, tt.store.student)
				}
			}
		})
	}
}
