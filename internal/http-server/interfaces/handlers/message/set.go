package message

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)



type SetRequest struct {
	ID    string `json:"id"`
}

type SetResponse struct {
	err int `json:"error, omitempty"`
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type messageSetter interface {
	Insert(ctx context.Context, text string) error
}

type Producer interface {
	ProduceMessage(text string) error
}

func NewSetter(userGetter messageSetter, produce Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := userGetter.Insert(context.Background(), req.ID)
		if err != nil {
			log.Println("Faled to insert data %v", err)
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		err = produce.ProduceMessage(req.ID)
		if err != nil {
			log.Println("Failed to send in kafka %v", err)
			http.Error(w, "Failed to send in kafka %v", http.StatusInternalServerError)
		}

		res := SetResponse{
			ID:      req.ID,
			Message: "Data inserted successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
