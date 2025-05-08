package delivery

import (
	"encoding/json"
	"net/http"
	"task1/response"
)

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.ErrorResponse{
		Error: message,
	})
	
}
