package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
	// return json.NewEncoder(w).Encode(data) // encode data to JSON and write to responses
}


func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}


func ValidationError(errs validator.ValidationErrors) Response {
	var erroMsges []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			erroMsges = append(erroMsges, fmt.Sprintf("%s is required", err.Field()))
		case "gte":
			erroMsges = append(erroMsges, fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param()))
		case "lte":
			erroMsges = append(erroMsges, fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param()))
		default:
			erroMsges = append(erroMsges, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(erroMsges, "; "),
	}
}
