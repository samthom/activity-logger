package errors

import "net/http"

// InternalError - emits internal server error
func InternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "Something went wrong."}`))
}
