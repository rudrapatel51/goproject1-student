package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rudrapatel51/goproject1-student/internal/storage"
	"github.com/rudrapatel51/goproject1-student/internal/types"
	"github.com/rudrapatel51/goproject1-student/internal/utils/response"
)

// New returns an http.HandlerFunc for creating a student
func New(studentStorage storage.Storage) http.HandlerFunc { // dependency injection of storage layer
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

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

		id, err := studentStorage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			if errors.Is(err, storage.ErrStudentEmailAlreadyExists) {
				response.WriteJson(w, http.StatusConflict, response.GeneralError(errors.New("student with this email already exists")))
				return
			}
			slog.Error("Failed to create student", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		student.ID = int(id)

		slog.Info("Student created", slog.Any("student", student))
		response.WriteJson(w, http.StatusCreated, student)
	}
}

func GetById(studentStorage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student id")))
			return
		}

		student, err := studentStorage.GetStudentByID(id)
		if err != nil {
			if errors.Is(err, storage.ErrStudentNotFound) {
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(errors.New("student not found")))
				return
			}
			slog.Error("Failed to get student by id", slog.Int("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(errors.New("failed to fetch student")))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetAll(studentStorage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := studentStorage.GetAllStudents()
		if err != nil {
			slog.Error("Failed to get all students", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(errors.New("failed to fetch students")))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func UpdateById(studentStorage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Parse and validate ID from URL path
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student id")))
			return
		}

		// Step 2: Decode JSON body into student struct
		var student types.Student
		if err = json.NewDecoder(r.Body).Decode(&student); err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("empty request body")))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Step 3: Validate the decoded fields
		if err = validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// Step 4: Call storage to update
		if err = studentStorage.UpdateById(id, student.Name, student.Email, student.Age); err != nil {
			if errors.Is(err, storage.ErrStudentNotFound) {
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(errors.New("student not found")))
				return
			}
			if errors.Is(err, storage.ErrStudentEmailAlreadyExists) {
				response.WriteJson(w, http.StatusConflict, response.GeneralError(errors.New("student with this email already exists")))
				return
			}
			slog.Error("Failed to update student", slog.Int("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(errors.New("failed to update student")))
			return
		}

		// Step 5: Fetch the updated record from DB and return it
		updated, err := studentStorage.GetStudentByID(id)
		if err != nil {
			slog.Error("Failed to fetch updated student", slog.Int("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(errors.New("failed to fetch updated student")))
			return
		}

		slog.Info("Student updated", slog.Int("id", id))
		response.WriteJson(w, http.StatusOK, updated)
	}
}

func DeleteById(studentStorage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Parse and validate ID from URL path
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student id")))
			return
		}

		// Step 2: Call storage to delete
		if err = studentStorage.DeleteById(id); err != nil {
			if errors.Is(err, storage.ErrStudentNotFound) {
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(errors.New("student not found")))
				return
			}
			slog.Error("Failed to delete student", slog.Int("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(errors.New("failed to delete student")))
			return
		}

		slog.Info("Student deleted", slog.Int("id", id))
		response.WriteJson(w, http.StatusOK, response.Response{Status: response.StatusOk, Error: ""})
	}
}