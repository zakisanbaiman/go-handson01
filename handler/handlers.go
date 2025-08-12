package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rsp := ErrResponse{
			Message: "failed to marshal response",
			Details: []string{err.Error()},
		}
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			log.Printf("failed to write error response: %v", err)
		}
		return
	}
	w.WriteHeader(status)
	if _, err := w.Write(bodyBytes); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
