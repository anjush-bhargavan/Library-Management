package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)


var ctx = context.Background()
//Client variable can used to save key value pairs in redis
var Client *redis.Client

//InitRedis function initializes redis server
func InitRedis() {
	Client = redis.NewClient(&redis.Options{
        Addr:	  "localhost:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })

	_,err := Client.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to redis"+err.Error())
	}
}