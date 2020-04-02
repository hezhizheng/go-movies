package services

import (
	"go_movies/models"
	"go_movies/utils"
)

func AllCategoryDate() []utils.Categories {
	categories := models.AllCategory()

	var nav []utils.Categories
	utils.Json.Unmarshal([]byte(categories), &nav)

	return nav
}
