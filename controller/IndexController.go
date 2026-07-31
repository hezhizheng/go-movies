package controller

import (
	"github.com/julienschmidt/httprouter"
	"go_movies/services"
	"go_movies/utils"
	"go_movies/views/tmpl"
	"net/http"
	"strconv"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	show := make(map[string]interface{})

	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)

	// 首页板块按当前采集源的分类动态生成
	sections := make([]services.IndexSection, 0, len(Categories))
	for _, c := range Categories {
		sections = append(sections, services.IndexSection{
			Name:   c.Name,
			Link:   c.Link,
			Subs:   c.Sub,
			Movies: services.MovieListsRange("detail_links:id:"+c.TypeId, 0, 14),
		})
	}
	show["sections"] = sections

	tmpl.GoTpl.ExecuteTemplate(w, "index", show)
}

func Display(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	path := r.URL.Path
	cate := r.URL.Query().Get("cate")
	_start := r.URL.Query().Get("start")
	_stop := r.URL.Query().Get("stop")

	show := make(map[string]interface{})

	// 导航栏类目显示
	Categories := services.AllCategoryData()
	show["categories"] = Categories
	show["allCategories"] = getAllCategory(Categories)

	// 默认当前采集源的第一个主分类
	key := ""
	if len(Categories) > 0 {
		key = "detail_links:id:" + Categories[0].TypeId
	}

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

	cateIntId := 0
	cateStrId := ""

	if len(cate) > 0 {
		cateStrId = services.TransformCategoryId(cate)
		cateIntId, _ = strconv.Atoi(cateStrId)
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

	// 根据不同类别显示不同 筛选类别项
	if subCate := currentGroupSubCategories(cateIntId, Categories); len(subCate) > 0 {
		show["currentSubCate"] = subCate
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

// 判断当前类目 id 属于哪个主分类，返回该主分类的子分类
func currentGroupSubCategories(cateIntId int, Categories []utils.Categories) []utils.Categories {
	for _, c := range Categories {
		mainId, _ := strconv.Atoi(c.TypeId)
		if cateIntId == mainId {
			return c.Sub
		}
		for _, sub := range c.Sub {
			subId, _ := strconv.Atoi(sub.TypeId)
			if cateIntId == subId {
				return c.Sub
			}
		}
	}
	return nil
}
