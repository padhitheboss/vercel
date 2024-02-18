package main

import (
	"fmt"

	queuemodel "github.com/padhitheboss/kube-job-scheduler/pkg/queueHelper/model"
	kubejob "github.com/padhitheboss/kube-job-scheduler/pkg/services/kubeJob"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println(err)
	// }
	queuemodel.StartQueue()
	reqFromQueue, _ := queuemodel.CollectRequest()
	fmt.Println(reqFromQueue)
	var k kubejob.Kubejob
	k.Initialize()
	k.CreateJobTemplate(reqFromQueue)
	k.RunJob()
}
