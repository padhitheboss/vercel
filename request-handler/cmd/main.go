package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/padhitheboss/request-handler/pkg/handler"
	queuemodel "github.com/padhitheboss/request-handler/pkg/queueHelper/model"
)

func main() {
	// godotenv.Load()
	queuemodel.StartQueue()
	handler.ConnectBlob()
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(handler.ServeHTTP)
	fmt.Println("Server listening on port 3001")
	err := http.ListenAndServe(":3001", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
