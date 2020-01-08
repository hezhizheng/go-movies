package utils

import (
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

// 声明一个全局的redisdb变量
var RedisDB *redis.Client

func InitRedisDB() (err error) {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     viper.GetString(`redis.addr`) + ":" + viper.GetString(`redis.port`),
		Password: viper.GetString(`redis.password`), // no password set
		DB:       viper.GetInt(`redis.db`),          // use default DB
	})

	_, err = RedisDB.Ping().Result()
	if err != nil {
		panic(err)
	}
	return nil
}

func CloseRedisDB() error {
	return RedisDB.Close()
}
