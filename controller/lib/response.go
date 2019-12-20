package lib

import "net/http"

// JSONContent - Helper for making the response to application/json format
func JSONContent(w http.ResponseWriter) http.ResponseWriter {
	return w
}
