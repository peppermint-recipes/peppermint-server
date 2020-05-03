package utils

import (
	"encoding/json"
	"net/http"
)

type ReturnMessage struct {
	Status     string
	Message    string
	StatusCode int
}

func PrepareReturn() map[string]interface{} {
	return map[string]interface{}{}
}

func ErrorMessage(message string) map[string]interface{} {
	return map[string]interface{}{"error": message}
}

func RespondWithError(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(data)
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)

	json.NewEncoder(w).Encode(data)
}
