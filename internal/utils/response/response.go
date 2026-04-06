package response

import (
	// "encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// return json.NewEncoder(w).Encode(data) // encode data to JSON and write to responses
}
