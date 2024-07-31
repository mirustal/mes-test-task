package message

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"kafka-app/internal/domain"
)

type GetRequest struct {
	ID string `json:"id"`
	Limit int `json:"limit"`
}

type GetResponse struct {
	Err int `json:"error"`
	Messages []domain.Message `json:"id"`
}

type UserGetter interface {
	GetUser(ctx context.Context, name string, limit int) ([]domain.Message, error)
}

func NewGetter(userGetter UserGetter) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		var req GetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Failed to get data", http.StatusInternalServerError)
			return
		}
	
	messages, err := userGetter.GetUser(context.Background(), req.ID, req.Limit)
	if err != nil{
		log.Println("Faled to insert data %v", err)
		http.Error(w, "failed to get data", http.StatusNotFound)
		return
	}

	res := GetResponse{
		Err: http.StatusAccepted,
		Messages: messages,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed encode response", http.StatusInternalServerError)
		return
	}
}
}