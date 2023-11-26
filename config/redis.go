package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)


var ctx = context.Background()
//Client variable can used to save key value pairs in redis
var Client *redis.Client

//InitRedis function initializes redis server
func InitRedis() {
	var err error
	MaxRetries:=5
	RetryDelay := time.Second * 5
	for i:=0;i<MaxRetries;i++{
		Client = redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:	  "localhost:6379",
			Password: "", // no password set
			DB:		  0,  // use default DB
		})

		_,err = Client.Ping(ctx).Result()
		if err == nil {
			break
		}
	
		fmt.Printf("Failed to connect to Redis (Attempt %d/%d): %s\n", i+1, MaxRetries, err.Error())
		time.Sleep(RetryDelay)
	}
	if err != nil {
		panic("Failed to connect to Redis after multiple attempts: " + err.Error())
	}
}