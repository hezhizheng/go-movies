package controller

import (
	"github.com/julienschmidt/httprouter"
	"go_movies/services"
	"go_movies/utils"
	"go_movies/utils/spider/tian_kong"
	"go_movies/views/tmpl"
	"net/http"
	"strconv"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	show := make(map[string]interface{})

	NewFilmKey := "detail_links:id:1"
	NewTVKey := "detail_links:id:3"
	NewCartoonKey := "detail_links:id:24"
	NewFilm := services.MovieListsRange(NewFilmKey, 0, 14)
	NewTV := services.MovieListsRange(NewTVKey, 0, 14)
	NewCartoon := services.MovieListsRange(NewCartoonKey, 0, 14)

	show["newFilm"] = NewFilm
	show["newTv"] = NewTV
	show["newCartoon"] = NewCartoon

	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)

	show["navFilm"] = getAssignTypeSubCategories(Categories, "film")
	show["navTv"] = getAssignTypeSubCategories(Categories, "tv")

	tmpl.GoTpl.ExecuteTemplate(w, "index", show)
}

func Display(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	path := r.URL.Path
	cate := r.URL.Query().Get("cate")
	_start := r.URL.Query().Get("start")
	_stop := r.URL.Query().Get("stop")

	show := make(map[string]interface{})

	key := "detail_links:id:5" // 默认首页

	start := int64(0)
	stop := int64(41)

	if len(_start) > 0 {
		StartInt64, _ := strconv.ParseInt(_start, 10, 64)
		start = StartInt64
	}

	if len(_stop) > 0 {
		StopInt64, _ := strconv.ParseInt(_stop, 10, 64)
		stop = StopInt64
	}

	prev := path + "?start=" + strconv.FormatInt(start-42, 10) + "&stop=" + strconv.FormatInt(stop-42, 10)
	next := path + "?start=" + strconv.FormatInt(start+42, 10) + "&stop=" + strconv.FormatInt(stop+42, 10)

	prevStatus := "1"
	nextStatus := "1"

	cateStrId := services.TransformCategoryId(cate)
	cateIntId, _ := strconv.Atoi(cateStrId)

	if len(cate) > 0 {
		key = "detail_links:id:" + cateStrId
		prev = path + "?cate=" + cate + "&start=" + strconv.FormatInt(start-42, 10) + "&stop=" + strconv.FormatInt(stop-42, 10)
		next = path + "?cate=" + cate + "&start=" + strconv.FormatInt(start+42, 10) + "&stop=" + strconv.FormatInt(stop+42, 10)
	}

	if start > stop || stop-start > 42 || start < 0 {
		start = 0
		stop = 41
	}

	MovieLists := services.MovieListsRange(key, start, stop)

	LenMovieLists := len(MovieLists)

	if start-42 < 0 {
		prevStatus = "0"
	}

	if LenMovieLists < 42 || LenMovieLists == 0 {
		nextStatus = "0"
	}

	show["movieLists"] = MovieLists
	show["prev"] = prev
	show["next"] = next
	show["prev_status"] = prevStatus
	show["next_status"] = nextStatus

	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)

	// 根据不同类别显示不同 筛选类别项
	if utils.InArray(cateIntId, tian_kong.GetAssignCategoryIds("film")) || cateIntId == 1 {
		show["currentSubCate"] = getAssignTypeSubCategories(Categories, "film")
	}

	if utils.InArray(cateIntId, tian_kong.GetAssignCategoryIds("tv")) || cateIntId == 2 {
		show["currentSubCate"] = getAssignTypeSubCategories(Categories, "tv")
	}

	if utils.InArray(cateIntId, tian_kong.GetAssignCategoryIds("cartoon")) || cateIntId == 4 {
		show["currentSubCate"] = getAssignTypeSubCategories(Categories, "cartoon")
	}

	tmpl.GoTpl.ExecuteTemplate(w, "display", show)
}

func Movie(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	link := r.URL.Query().Get("link")
	if link == "" {
		w.WriteHeader(404)
		w.Write([]byte("404"))
		return
	}

	show := make(map[string]interface{})
	MovieDetail := services.MovieDetail(link)

	if len(MovieDetail["info"].(map[string]string)) == 0 {
		w.WriteHeader(404)
		tmpl.GoTpl.ExecuteTemplate(w, "404", show)
		return
	}

	show["MovieDetail"] = MovieDetail
	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)
	tmpl.GoTpl.ExecuteTemplate(w, "detail", show)
}

func Play(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	show := make(map[string]interface{})
	PlayUrl := r.URL.Query().Get("play_url")
	PlayType := "kuyun"
	if strings.Contains(PlayUrl, ".mp4") {
		PlayType = "mp4"
	} else if strings.Contains(PlayUrl, ".m3u8") {
		PlayType = "m3u8"
	}

	show["playUrl"] = PlayUrl
	show["playType"] = PlayType

	link := r.URL.Query().Get("link")
	episode := r.URL.Query().Get("episode")
	MovieDetail := services.MovieDetail(link)

	if len(MovieDetail["info"].(map[string]string)) == 0 {
		w.WriteHeader(404)
		tmpl.GoTpl.ExecuteTemplate(w, "404", show)
		return
	}

	show["MovieDetail"] = MovieDetail
	show["episode"] = episode
	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)

	tmpl.GoTpl.ExecuteTemplate(w, "play", show)
}

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	show := make(map[string]interface{})
	q := r.URL.Query().Get("q")

	var MovieLists []services.MovieListStruct
	if strings.TrimSpace(q) != "" {
		MovieLists = services.SearchMovies(q)
	}

	show["movieLists"] = MovieLists
	show["q"] = q
	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)
	tmpl.GoTpl.ExecuteTemplate(w, "search", show)
}

func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	show := make(map[string]interface{})
	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)
	tmpl.GoTpl.ExecuteTemplate(w, "about", show)
}

func getAllCategory(Categories []utils.Categories) []utils.Categories {
	var allC []utils.Categories
	for _, c := range Categories {
		allC = append(allC, c)
		for _, subC := range c.Sub {
			allC = append(allC, subC)
		}
	}
	return allC
}

func getAssignTypeSubCategories(Categories []utils.Categories, _type string) []utils.Categories {
	var cate []utils.Categories
	switch _type {
	case "film":
		cate = Categories[0].Sub
	case "tv":
		cate = Categories[1].Sub
	case "cartoon":
		cate = Categories[2].Sub
	}
	return cate
}
