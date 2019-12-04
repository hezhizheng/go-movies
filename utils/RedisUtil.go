package utils

import "github.com/go-redis/redis/v7"

// 声明一个全局的redisdb变量
var RedisDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       10, // use default DB
})
