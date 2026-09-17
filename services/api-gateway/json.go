package main

import (
	"encoding/json"
	"net/http"
)

func writeJSONResponse(w http.ResponseWriter, status int, response any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(response)
}
