package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/rudrapatel51/goproject1-student/internal/utils/response"
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
			slog.Error("Failed to decode request body", slog.String("error", err.Error()))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadGateway, err.Error())
			return
		}

		slog.Info("Student created", slog.Any("student", student))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(student)
	}
}
