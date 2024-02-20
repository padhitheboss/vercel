package queueconfig

import (
	"log"

	"github.com/padhitheboss/code-builder/pkg/service/queue/redisQueue"
)

type QueueFunction interface {
	Connect()
	ReadFromQueue() (string, error)
	WriteToQueue(string) error
	UpdateDB(string, string) error
	GetFromDB(string) (string, error)
}

// Initialize Queue

func InitializeQueue(queueType string) QueueFunction {
	switch queueType {
	case "redis":
		var redisQ redisQueue.RedisConfig = redisQueue.CreateConfig()
		var q QueueFunction = &redisQ
		q.Connect()
		return q
	default:
		log.Panicf("invalid queue type: %v", queueType)
		return nil
	}
	// return nil, fmt.Errorf("unable to initialize queue")
}
