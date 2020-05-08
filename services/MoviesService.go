package services

import (
	"go_movies/models"
	"go_movies/utils"
	"go_movies/utils/spider"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

type MovieListStruct struct {
	Link      string `json:"link"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Cover     string `json:"cover"`
	UpdatedAt string `json:"updated_at"`
	Starring  string `json:"starring"`
}

type SortMovieListStruct []MovieListStruct

// 获取此 slice 的长度 @deprecated
func (p SortMovieListStruct) Len() int { return len(p) }

// 根据元素的年龄降序排序 （此处按照自己的业务逻辑写）@deprecated
func (p SortMovieListStruct) Less(i, j int) bool {

	// 模板时间
	timeTemplate := "2006-01-02"

	stamp1, _ := time.ParseInLocation(timeTemplate, p[i].UpdatedAt, time.Local)
	stamp2, _ := time.ParseInLocation(timeTemplate, p[j].UpdatedAt, time.Local)

	return stamp1.Unix() > stamp2.Unix()
}

// 交换数据 @deprecated
func (p SortMovieListStruct) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

var (
	mutex sync.Mutex
)

// @deprecated
func MovieLists(key string) []MovieListStruct {
	//mutex.Lock() //对共享资源加锁
	//defer mutex.Unlock()

	var data []MovieListStruct
	var NewData []MovieListStruct

	//movieKeys := models.FindMoviesKey(key)
	movieKeys := RangeSCanMoviesKey(key)

	for _, val := range movieKeys {
		movieKey := models.FindMoviesStringValue(val) // json

		var movieKeyMap MovieListStruct
		err := utils.Json.Unmarshal([]byte(movieKey), &movieKeyMap)
		if err != nil {
			log.Println("movieKeyMap err: ", err)
		}

		mutex.Lock()
		data = append(data, movieKeyMap)
		mutex.Unlock()
	}

	for _, dataMapVal := range data {
		link := dataMapVal.Link

		detail := models.FindMoviesHashValue("movies_detail:" + link)

		//if detail["name"] == "" {
		//	go utils.MoviesInfo(link) // 重新采集
		//}

		dataMapVal.Cover = detail["cover"] // 重新赋值

		mutex.Lock()
		NewData = append(NewData, dataMapVal)
		mutex.Unlock()
	}

	// golang 处理排序也.....
	sort.Sort(SortMovieListStruct(NewData))

	return NewData
}

// 查询指定范围的数据
func MovieListsRange(key string, start, stop int64) []MovieListStruct {
	var data []MovieListStruct

	var movieKeyMap MovieListStruct

	sStart := strconv.FormatInt(start, 10)
	sStop := strconv.FormatInt(stop, 10)

	cacheKey := "movie_lists_key:" + key + ":start:" + sStart + ":stop:" + sStop

	movieList := models.FindMoviesStringValue(cacheKey)

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

		byteData, _ := utils.Json.MarshalIndent(data, "", " ")
		models.SaveMovies(cacheKey, string(byteData))
	}

	return data
}

func TransformCategoryId(link string) string {
	return utils.TransformId(link)
}

func MovieDetail(link string) map[string]interface{} {
	mutex.Lock() //对共享资源加锁
	defer mutex.Unlock()
	data := make(map[string]interface{})

	details := models.FindMoviesKey("movies_detail:" + link + "*")

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

// SCan 代替 keys @deprecated
func RangeSCanMoviesKey(key string) []string {
	var (
		all []string
		i   uint64
	)

	i = 0
	for {
		s, c, _ := models.SCanMoviesKey(i, key+"*", 1000)

		log.Println(s, c)
		// 如果只有一个，停止循环
		if len(s) == 1 {
			all = s
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

// 搜索影片
func SearchMovies(key string) []MovieListStruct {

	var data []MovieListStruct

	var movieKeyMap MovieListStruct

	movieKeys := models.FindMoviesKey("*" + ":movie_name:" + key + "*")

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
