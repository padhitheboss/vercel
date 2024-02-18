package main

import (
	"net/http"

	"example.com/uploader/pkg/handler"
	queuemodel "example.com/uploader/pkg/queueHelper/model"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// var wg sync.WaitGroup
	// rq := redisQueue.CreateConfig()
	// rq.Connect()
	// wg.Add(1)
	// go func() {
	// s, _ := rq.ReadFromQueue()
	// fmt.Println(s)
	// time.Sleep(10 * time.Second)
	// s, _ = rq.ReadFromQueue()
	// fmt.Println(s)
	// rq.WriteToQueue("hello i am writing")
	// rq.WriteToQueue("r u there")
	// rq.WriteToQueue("hello")
	// }()
	// wg.Wait()
	// azureBlob.CreateZipFile("/workspaces/codespaces-blank/test.zip", "/workspaces/codespaces-blank/uploader/")
	queuemodel.StartQueue()
	http.HandleFunc("/deploy", handler.RequestDeploy)
	http.ListenAndServe(":8080", nil)
	// fmt.Println(queuemodel.CreateResponse("1", "success", "abcd-efg-hij"))
}
