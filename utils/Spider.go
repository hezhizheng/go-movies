package utils

import "strings"

// redis key
// 分类key
const CategoriesKey = "categories"

type Categories struct {
	Link   string       `json:"link"`
	Name   string       `json:"name"`
	TypeId string       `json:"type_id"`
	Sub    []Categories `json:"sub"`
}

// 获取url中的链接
func TransformId(Url string) string {
	UrlStrSplit := strings.Split(Url, "-id-")[1]
	return strings.TrimRight(UrlStrSplit, ".html")
}