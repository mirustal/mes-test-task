package handlers

import (
	"fmt"
	"net/http"
)

func SetMessageHandler(w http.ResponseWriter, r *http.Request)  {
    fmt.Println("hello")
}


