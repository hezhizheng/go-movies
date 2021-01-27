package utils

import (
	"strings"
)

func InArray(needle interface{}, hystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}

//
// /index.php/vod/detail/id/1405.html => /?m=vod-detail-id-16250.html
func Index2vod(Url string) string {

	UrlStrSplit := strings.Split(Url, "/id/")

	if len(UrlStrSplit) >= 2 {
		id := strings.TrimRight(UrlStrSplit[1], ".html")

		return "/?m=vod-detail-id-"+id+".html"
	}

	return Url

}

func Index2vodList(Url string) string {
	UrlStrSplit := strings.Split(Url, "/id/")

	if len(UrlStrSplit) >= 2 {
		id := strings.TrimRight(UrlStrSplit[1], ".html")

		return "/?m=vod-type-id-"+id+".html"
	}
	return Url
}

// /?m=vod-detail-id-16250.html => /index.php/vod/detail/id/1405.html
func Vod2index(Url string) string {
	UrlStrSplit := strings.Split(Url, "-id-")

	if len(UrlStrSplit) >= 2 {
		id := strings.TrimRight(UrlStrSplit[1], ".html")
		return "/index.php/vod/detail/id/"+id+".html"
	}

	return Url
}

func Vod2indexList(Url string) string {
	UrlStrSplit := strings.Split(Url, "-id-")

	if len(UrlStrSplit) >= 2 {
		id := strings.TrimRight(UrlStrSplit[1], ".html")
		return "/index.php/vod/type/id/"+id+".html"
	}

	return Url
}