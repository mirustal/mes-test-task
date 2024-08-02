package routes

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

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
	// mux.HandleFunc("/read", message.NewMark(db, consumer))

	mux.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/swagger.json"),
	))


	handler := middlewares.EnableCORS(mux)
	handler = middlewares.LoggingRequest(handler)

	return handler
}
