package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"time"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	time.Sleep(9 * time.Second)

	//TODO: call the trip service via grpc later
	var reqBody previewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse json data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	jsonBody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

	resp, err := http.Post("http://localhost:8083/preview", "application/json", reader)
	if err != nil {
		log.Print(err)
		return
	}

	defer resp.Body.Close()

	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, "failed to parse json data", http.StatusBadRequest)
	}

	response := contracts.APIResponse{
		// Error: &contracts.APIError{
		// 	Code: hex(http.StatusCreated),
		// },
		Data: "ok",
	}
	writeJSONResponse(w, http.StatusCreated, response)
}
