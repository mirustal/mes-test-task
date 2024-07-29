package routes

import (
	"net/http"

	"kafka-app/internal/http-server/interfaces/handlers"
	"kafka-app/internal/http-server/interfaces/middlewares"
)

type Handler struct {
	handler *http.ServeMux
}

func InitRoutes() http.Handler {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/get", handlers.GetUserHandler)
	mux.HandleFunc("/set", handlers.SetMessageHandler)

	return middlewares.LoggingRequest(mux)
}


