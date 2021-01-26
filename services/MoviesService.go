package services

import (
	"github.com/spf13/viper"
	"go_movies/models"
	"go_movies/utils"
	"go_movies/utils/spider"
	"strconv"
	"strings"
	"sync"
)

const paginateCacheKey  = "paginate"

type MovieListStruct struct {
	Link      string `json:"link"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Cover     string `json:"cover"`
	UpdatedAt string `json:"updated_at"`
	Starring  string `json:"starring"`
}

type SortMovieListStruct []MovieListStruct

var (
	mutex sync.Mutex
)

// 查询指定范围的数据
func MovieListsRange(key string, start, stop int64) []MovieListStruct {
	var data []MovieListStruct

	var movieKeyMap MovieListStruct

	sStart := strconv.FormatInt(start, 10)
	sStop := strconv.FormatInt(stop, 10)

	// field
	cacheKey := "movie_lists_key:" + key + ":start:" + sStart + ":stop:" + sStop
	cacheHashKey := paginateCacheKey

	movieList := models.FindMoviesHashFieldValue(cacheHashKey,cacheKey)

	if movieList != "" {
		utils.Json.Unmarshal([]byte(movieList), &data)
		return data
	}

	movieKeys, _ := models.ZREVRANGEMoviesKey(key, start, stop)

	for _, val := range movieKeys {

		MovieDetail := MovieDetail(val)

		info := MovieDetail["info"].(map[string]string)
		detail := MovieDetail["detail"].(map[string]interface{})

		movieKeyMap.Name = info["name"]
		movieKeyMap.Link = info["link"]
		movieKeyMap.Cover = info["cover"]
		if detail["update"] != nil {
			movieKeyMap.UpdatedAt = detail["update"].(string)
		}
		if detail["type"] != nil {
			movieKeyMap.Category = detail["type"].(string)
		}
		if detail["starring"] != nil {
			movieKeyMap.Starring = detail["starring"].(string)
		}

		mutex.Lock()
		data = append(data, movieKeyMap)
		mutex.Unlock()

	}

	// 这个应该是在for外面调用才对，之前居然也不报错。。。。。。
	byteData, _ := utils.Json.MarshalIndent(data, "", " ")
	models.SaveMoviesHash(cacheHashKey, cacheKey ,string(byteData))

	return data
}

func TransformCategoryId(link string) string {
	return utils.TransformId(link)
}

func MovieDetail(link string) map[string]interface{} {
	mutex.Lock() //对共享资源加锁
	defer mutex.Unlock()
	data := make(map[string]interface{})

	//details := models.FindMoviesKey("movies_detail:" + link + "*")
	details := models.RangeSCanMoviesKey("movies_detail:" + link + "*")

	detail := make(map[string]string)
	if len(details) > 0 {
		detail = models.FindMoviesHashValue(details[0])
	}

	if detail["name"] == "" {
		// 重新采集
		detailId := TransformCategoryId(link)
		spider.Create().PageDetail(detailId)
	}

	var kuYunMap []map[string]interface{}
	utils.Json.Unmarshal([]byte(detail["kuyun"]), &kuYunMap)

	var ckm3u8Map []map[string]interface{}
	utils.Json.Unmarshal([]byte(detail["ckm3u8"]), &ckm3u8Map)

	var downloadMap []map[string]interface{}
	utils.Json.Unmarshal([]byte(detail["download"]), &downloadMap)

	var _detailMap map[string]interface{}
	utils.Json.Unmarshal([]byte(detail["detail"]), &_detailMap)

	data["kuyun"] = kuYunMap
	data["ckm3u8"] = ckm3u8Map
	data["download"] = downloadMap
	data["detail"] = _detailMap

	delete(detail, "kuyun")
	delete(detail, "ckm3u8")
	delete(detail, "download")
	delete(detail, "detail")

	data["info"] = detail

	isFilm := "1"

	if len(_detailMap) > 0 {
		if strings.Index(_detailMap["type"].(string), "片") == -1 { // ...片  or  ... 剧
			isFilm = "0"
		}
	}

	data["is_film"] = isFilm

	return data
}

// 搜索影片
func SearchMovies(key string) []MovieListStruct {

	var data []MovieListStruct

	var movieKeyMap MovieListStruct

	//movieKeys := models.FindMoviesKey("*" + ":movie_name:" + key + "*")
	movieKeys := models.RangeSCanMoviesKey("*" + ":movie_name:" + key + "*")

	for _, val := range movieKeys {

		MovieDetail := MovieDetail(TransformLink(val))

		info := MovieDetail["info"].(map[string]string)
		detail := MovieDetail["detail"].(map[string]interface{})

		movieKeyMap.Name = info["name"]
		movieKeyMap.Link = info["link"]
		movieKeyMap.Cover = info["cover"]
		if detail["update"] != nil {
			movieKeyMap.UpdatedAt = detail["update"].(string)
		}
		if detail["type"] != nil {
			movieKeyMap.Category = detail["type"].(string)
		}
		if detail["starring"] != nil {
			movieKeyMap.Starring = detail["starring"].(string)
		}

		mutex.Lock()
		data = append(data, movieKeyMap)
		mutex.Unlock()
	}

	return data

}

// 获取实际链接url
func TransformLink(Url string) string {
	UrlStrSplit := strings.Split(Url, "movies_detail:")[1]

	return strings.Split(UrlStrSplit, ":movie_name:")[0]
}

func MoviesRecommend() interface{} {
	recommend := viper.Get(`recommend`)
	c := new([]interface{})
	if recommend == nil {
		return *c
	}
	return recommend
}