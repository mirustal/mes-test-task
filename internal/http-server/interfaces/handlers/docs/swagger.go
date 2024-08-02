package docs

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)


func SwaggerHandler(w http.ResponseWriter, r *http.Request) {
    swaggerPath := filepath.Join("docs", "index.html")
    swaggerFile, err := ioutil.ReadFile(swaggerPath)
    if err != nil {
        http.Error(w, "Could not read swagger.docs file", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/html")
    w.Write(swaggerFile)
}