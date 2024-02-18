package redisQueue

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func CreateConfig() RedisConfig {
	dbId, err := strconv.Atoi(os.Getenv("Q_DB"))
	if err != nil {
		log.Fatalf("invalid db check config")
	}
	return RedisConfig{
		addr:         os.Getenv("Q_ADDRESS"),
		port:         os.Getenv("Q_PORT"),
		password:     os.Getenv("Q_PASSWORD"),
		db:           int(dbId),
		readChannel:  os.Getenv("Q_R_CHANNEL"),
		writeChannel: os.Getenv("Q_W_CHANNEL"),
	}
}

func (r *RedisConfig) Connect() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.addr + ":" + r.port, // Redis server address
		Password: r.password,            // Redis server password
		DB:       r.db,                  // Redis server DB
	})
	ctx := context.Background()
	r.pubsub = r.client.Subscribe(ctx, r.readChannel)
	// Wait for confirmation that subscription is created
	_, err := r.pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("connected to redis")
}

func (r *RedisConfig) ReadFromQueue() (string, error) {
	message, err := r.pubsub.ReceiveMessage(context.Background())
	if err == nil {
		return message.Payload, nil
	}
	return "", fmt.Errorf("channel closed unexpectedly")
}

func (r *RedisConfig) WriteToQueue(message string) error {
	res := r.client.Publish(context.Background(), r.writeChannel, message)
	return res.Err()
}
