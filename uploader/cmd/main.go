package main

import (
	"fmt"
	"net/http"

	"example.com/uploader/pkg/handler"
	queuemodel "example.com/uploader/pkg/queueHelper/model"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	queuemodel.StartQueue()
	http.HandleFunc("/deploy", handler.RequestDeploy)
	http.HandleFunc("/status", handler.DeployStatus)
	fmt.Println("starting http server on port: 8080")
	http.ListenAndServe(":8080", nil)
}
