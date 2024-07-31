package routes

import (
	"net/http"

	"kafka-app/internal/adapters/db/postgres"
	"kafka-app/internal/adapters/kafka/producer"
	"kafka-app/internal/http-server/interfaces/handlers/message"
	"kafka-app/internal/http-server/interfaces/middlewares"
)

type Handler struct {
	handler *http.ServeMux
}

func InitRoutes(db *postgres.PostgresMessageRep, producer *producer.Producer) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/get", message.NewGetter(db))
	mux.HandleFunc("/set", message.NewSetter(db, producer))

	return middlewares.LoggingRequest(mux)
}
