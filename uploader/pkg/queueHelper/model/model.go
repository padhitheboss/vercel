package queuemodel

import (
	"encoding/json"
	"strings"

	"example.com/uploader/pkg/model"
	queueconfig "example.com/uploader/pkg/queueHelper/config"
)

var q queueconfig.QueueFunction

func StartQueue() {
	q = queueconfig.InitializeQueue("redis")
}

type Response struct {
	RequestId   string `json:"request_id"`
	Status      string `json:"status"`
	BlobUrl     string `json:"blob_url"`
	ResourceId  string `json:"resource_id"`
	ProjectType string `json:"project_type"`
}

func CreateResponse(res model.Request, uploadPath, status string) Response {
	return Response{
		RequestId:   res.RequestId,
		Status:      status,
		BlobUrl:     strings.TrimPrefix(uploadPath, ""),
		ResourceId:  res.RequestId,
		ProjectType: res.ProjectType,
	}
}

func SendResponse(res Response) error {
	responseString, _ := json.Marshal(res)
	return q.WriteToQueue(string(responseString))
}
