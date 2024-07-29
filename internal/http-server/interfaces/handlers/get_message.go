package handlers

import (
	"fmt"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request)  {
    fmt.Println("hello")
}


