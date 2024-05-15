package utils

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func checkRedisConnectionStatus(client *redis.Client) error {
	err := client.Conn().Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	return nil
}

func ConnectRedis() (*redis.Client, error) {

	fmt.Println("---------------- redis ---------------------")
	var err error

	var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		// Addr:     "host.docker.internal:6379", 
		Password: "",
		DB:       0, // use default db
	})

	err = checkRedisConnectionStatus(client)
	if err != nil {
		fmt.Println("Connection Error : ", err)
		return client, err
	}

	return client, nil
}

