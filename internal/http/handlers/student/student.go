package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/rudrapatel51/goproject1-student/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// New returns an http.HandlerFunc for creating a student
func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student Student

		err := json.NewDecoder(r.Body).Decode(&student) // decode request body into student struct
		if err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("empty request body")))
				return
			}
			slog.Error("Failed to decode request body", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}



		slog.Info("Student created", slog.Any("student", student))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(student)
	}
}
