package queuemodel

import (
	"encoding/json"
	"fmt"
	"strings"

	"example.com/uploader/pkg/model"
	queueconfig "example.com/uploader/pkg/queueHelper/config"
)

var q queueconfig.QueueFunction

func StartQueue() {
	q = queueconfig.InitializeQueue("redis")
}

type Response struct {
	RequestId    string `json:"request_id"`
	Status       string `json:"status"`
	BlobUrl      string `json:"blob_url"`
	ResourceId   string `json:"resource_id"`
	ProjectType  string `json:"project_type"`
	EndpointName string `json:"endpoint_name"`
}

func CreateResponse(res model.Request, uploadPath, status string) Response {
	return Response{
		RequestId:    res.RequestId,
		Status:       status,
		BlobUrl:      strings.TrimPrefix(uploadPath, ""),
		ResourceId:   res.RequestId,
		ProjectType:  res.ProjectType,
		EndpointName: res.EndPointName,
	}
}

func SendResponse(res Response) error {
	responseString, _ := json.Marshal(res)
	return q.WriteToQueue(string(responseString))
}

func UpdateDB(l model.LogDB) error {
	value, err := json.Marshal(l)
	if err != nil {
		fmt.Println("error converting json to string:", err)
		return err
	}
	if l.Status == "" && l.LayerLogs[len(l.LayerLogs)-1].Status == "failed" {
		l.Status = "failed"
	}
	return q.UpdateDB(l.EndpointName, string(value))
}

func GetFromDB(id string) (model.LogDB, error) {
	var l model.LogDB
	value, err := q.GetFromDB(id)
	if err != nil {
		fmt.Println("error converting string to json:", err)
		return l, err
	}
	json.Unmarshal([]byte(value), &l)
	return l, nil
}
