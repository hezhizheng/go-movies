package models

import (
	"go_movies/utils"
)

func AllCategory() string {
	return utils.RedisDB.Get(utils.CategoriesKey).Val()
}
