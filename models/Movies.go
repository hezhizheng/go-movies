package models

import (
	"go_movies/utils"
	"sync"
)

func FindMoviesKey(key string) []string {
	return utils.RedisDB.Keys(key).Val()
}

func SCanMoviesKey(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return utils.RedisDB.Scan(cursor, match, count).Result()
}

func ZREVRANGEMoviesKey(key string, start, stop int64) ([]string, error) {
	return utils.RedisDB.ZRevRange(key, start, stop).Result()
}

func FindMoviesStringValue(key string) string {
	return utils.RedisDB.Get(key).Val()
}

func FindMoviesHashValue(key string) map[string]string {
	return utils.RedisDB.HGetAll(key).Val()
}

func SaveMovies(key string, value string) error {
	return utils.RedisDB.SetNX(key, value, 0).Err()
}

func SaveMoviesHash(key, field string, value interface{}) error {
	return utils.RedisDB.HSetNX(key, field, value).Err()
}

func FindMoviesHashFieldValue(key, field string) string {
	return utils.RedisDB.HGet(key,field).Val()
}

func RangeSCanMoviesKey(key string) []string {
	var (
		all []string
		i   uint64
		mutex sync.Mutex
	)

	i = 0
	for {
		s, c, _ := SCanMoviesKey(i, key, 1000)
		//log.Println("s c",s, c)
		// 游标为0，停止循环
		if c == 0 {
			for _, val := range s {
				mutex.Lock()
				all = append(all, val)
				mutex.Unlock()
			}
			break
		} else {
			i = c
			for _, val := range s {
				mutex.Lock()
				all = append(all, val)
				mutex.Unlock()
			}
		}
	}
	return all
}