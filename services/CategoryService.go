package services

import (
	"go_movies/models"
	"go_movies/utils"
)

func AllCategoryDate() []map[string]interface{} {
	categories := models.AllCategory()

	var nav []map[string]interface{}
	utils.Json.Unmarshal([]byte(categories), &nav)

	return nav
}
