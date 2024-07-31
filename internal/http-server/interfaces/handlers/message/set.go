package message

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)



type SetRequest struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
}

type SetResponse struct {
	err int `json:"error, omitempty"`
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type messageSetter interface {
	Insert(ctx context.Context, text string) (string, error)
}

type Producer interface {
	ProduceMessage(id, text string) error
}

func NewSetter(userGetter messageSetter, produce Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		id, err := userGetter.Insert(context.Background(), req.Text)
		if err != nil {
			log.Printf("Failed to insert data: %v", err)
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		err = produce.ProduceMessage(id, req.Text)
		if err != nil {
			log.Printf("Failed to send to Kafka: %v", err)
			http.Error(w, "Failed to send to Kafka", http.StatusInternalServerError)
			return
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
