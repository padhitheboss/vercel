package redisQueue

import "github.com/redis/go-redis/v9"

type RedisConfig struct {
	client       *redis.Client
	addr         string
	port         string
	password     string
	db           int
	readChannel  string
	writeChannel string
	pubsub       *redis.PubSub
}
