package logger

import (
	"time"
)

type LogDB struct {
	RequestId    string
	Status       string
	EndpointName string
	LayerLogs    []LayerLog
}

type LayerLog struct {
	Status    string
	Layer     string
	Rawlogs   string
	Timestamp string
}

func CreateLayerLog() LayerLog {
	return LayerLog{
		Layer:     "builder",
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func StartLog(reqId, EndpointName string) LogDB {
	return LogDB{
		RequestId:    reqId,
		EndpointName: EndpointName,
	}
}
