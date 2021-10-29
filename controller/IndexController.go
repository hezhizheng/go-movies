package controller

import (
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_movies/services"
	"go_movies/utils"
	"go_movies/utils/spider/tian_kong"
	heroTpl "go_movies/views/hero"
	"go_movies/views/tmpl"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// 首页
func Index1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	path := r.URL.Path
	cate := r.URL.Query()["cate"]
	_start := r.URL.Query()["start"]
	_stop := r.URL.Query()["stop"]

	// 需要展示的数据
	show := make(map[string]interface{})

	// 所有类别/导航
	Categories := services.AllCategoryData()

	key := "detail_links:id:14" // 默认首页

	start := int64(0)
	stop := int64(14)

	if len(_start) > 0 {
		StartInt64, _ := strconv.ParseInt(_start[0], 10, 64)
		start = StartInt64
	}

	if len(_stop) > 0 {
		StopInt64, _ := strconv.ParseInt(_stop[0], 10, 64)
		stop = StopInt64
	}

	prev := path + "?start=" + strconv.FormatInt(start-15, 10) + "&stop=" + strconv.FormatInt(stop-15, 10)
	next := path + "?start=" + strconv.FormatInt(start+15, 10) + "&stop=" + strconv.FormatInt(stop+15, 10)

	prevStatus := "1"
	nextStatus := "1"

	navLink := "/"
	if len(cate) > 0 {
		key = "detail_links:id:" + services.TransformCategoryId(cate[0])
		navLink = cate[0]
		prev = path + "?cate=" + cate[0] + "&start=" + strconv.FormatInt(start-15, 10) + "&stop=" + strconv.FormatInt(stop-15, 10)
		next = path + "?cate=" + cate[0] + "&start=" + strconv.FormatInt(start+15, 10) + "&stop=" + strconv.FormatInt(stop+15, 10)
	}

	if start > stop || stop-start > 15 || start < 0 {
		start = 0
		stop = 15
	}

	MovieLists := services.MovieListsRange(key, start, stop)

	LenMovieLists := len(MovieLists)

	if start-15 < 0 {
		prevStatus = "0"
	}

	if LenMovieLists < 15 || LenMovieLists == 0 {
		nextStatus = "0"
	}

	NewFilmKey := "detail_links:id:1"
	NewTVKey := "detail_links:id:2"

	NewFilm := services.MovieListsRange(NewFilmKey, 0, 49)
	NewTV := services.MovieListsRange(NewTVKey, 0, 49)
	recommend := services.MoviesRecommend()
	show["recommend"] = recommend

	show["categories"] = Categories
	show["page"] = "页面"
	show["movieLists"] = MovieLists
	show["new_film"] = NewFilm
	show["new_tv"] = NewTV
	show["prev"] = prev
	show["next"] = next
	show["prev_status"] = prevStatus
	show["next_status"] = nextStatus
	show["nav_link"] = navLink
	show["film_title"] = ""



	//GoTpl.ExecuteTemplate(os.Stdout, "index", "hzz go movies")
	//display ,rc:= tmpl.GoTpl.ParseFiles("./views/tmpl/index.html","./views/tmpl/nav.html")

	vv := tmpl.GoTpl.ExecuteTemplate(w,"index","hzz go movies33")
	log.Println("eeeeeeee",vv)

	//display.ExecuteTemplate(w,"nav","hzz go movies33")
	//tmpl.GoTpl.ExecuteTemplate(w,"nav","hzz go movies2")

	//tmpl.GoTpl.Execute(w,"hzz go movies2")

	//GoTpl.ParseFiles("index.html")
	//buffer := new(bytes.Buffer)
	//heroTpl.Index(show, buffer)
	//w.Write([]byte(`fdsfsdfs`))


}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// 需要展示的数据
	show := make(map[string]interface{})
	// 所有类别/导航
	Categories := services.AllCategoryData()

	NewFilmKey := "detail_links:id:1"
	NewTVKey := "detail_links:id:2"
	NewCartoonKey := "detail_links:id:4"

	NewFilm := services.MovieListsRange(NewFilmKey, 0, 14)
	NewTV := services.MovieListsRange(NewTVKey, 0, 14)
	NewCartoon := services.MovieListsRange(NewCartoonKey, 0, 14)

	show["categories"] = Categories

	show["newFilm"] = NewFilm
	show["newTv"] = NewTV
	show["newCartoon"] = NewCartoon


	// 合并所有的类目
	var allC []utils.Categories
	for _, c := range show["categories"].([]utils.Categories) {
		allC = append(allC,c)
		for _,subC := range c.Sub{
			allC = append(allC,subC)
		}
	}
	show["allCategories"] = allC
	show["navFilm"] = show["categories"].([]utils.Categories)[0].Sub
	show["navTv"] = show["categories"].([]utils.Categories)[1].Sub
	//show["navFilm"] = show["categories"].([]utils.Categories)[2].Sub
	//log.Println(show["categories"].([]utils.Categories)[0].Sub)


	tmpl.GoTpl.ExecuteTemplate(w,"index",show)


}

func Display(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	path := r.URL.Path
	cate := r.URL.Query()["cate"]
	_start := r.URL.Query()["start"]
	_stop := r.URL.Query()["stop"]

	// 需要展示的数据
	show := make(map[string]interface{})

	// 所有类别/导航
	Categories := services.AllCategoryData()

	key := "detail_links:id:14" // 默认首页

	start := int64(0)
	stop := int64(41)

	if len(_start) > 0 {
		StartInt64, _ := strconv.ParseInt(_start[0], 10, 64)
		start = StartInt64
	}

	if len(_stop) > 0 {
		StopInt64, _ := strconv.ParseInt(_stop[0], 10, 64)
		stop = StopInt64
	}

	prev := path + "?start=" + strconv.FormatInt(start-42, 10) + "&stop=" + strconv.FormatInt(stop-42, 10)
	next := path + "?start=" + strconv.FormatInt(start+42, 10) + "&stop=" + strconv.FormatInt(stop+42, 10)

	prevStatus := "1"
	nextStatus := "1"

	cateStrId := services.TransformCategoryId(cate[0])

	if len(cate) > 0 {
		key = "detail_links:id:" + cateStrId
		prev = path + "?cate=" + cate[0] + "&start=" + strconv.FormatInt(start-42, 10) + "&stop=" + strconv.FormatInt(stop-42, 10)
		next = path + "?cate=" + cate[0] + "&start=" + strconv.FormatInt(start+42, 10) + "&stop=" + strconv.FormatInt(stop+42, 10)
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


	show["categories"] = Categories
	show["movieLists"] = MovieLists
	show["prev"] = prev
	show["next"] = next
	show["prev_status"] = prevStatus
	show["next_status"] = nextStatus

	var allC []utils.Categories
	for _, c := range show["categories"].([]utils.Categories) {
		allC = append(allC,c)
		for _,subC := range c.Sub{
			allC = append(allC,subC)
		}
	}
	show["allCategories"] = allC


	if utils.InArray(cateStrId,tian_kong.GetAssignCategoryStrIds("film")) {
		show["currentSubCate"] = show["categories"].([]utils.Categories)[0].Sub
	}

	if utils.InArray(cateStrId,tian_kong.GetAssignCategoryStrIds("tv")) {
		show["currentSubCate"] = show["categories"].([]utils.Categories)[1].Sub
	}

	if utils.InArray(cateStrId,tian_kong.GetAssignCategoryStrIds("cartoon")) {
		show["currentSubCate"] = show["categories"].([]utils.Categories)[2].Sub
	}

	//log.Println(path,cate)


	tmpl.GoTpl.ExecuteTemplate(w,"display",show)

}

func Movie(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	_link := r.URL.Query()["link"]

	if len(_link) == 0 {
		fmt.Fprint(w, "404")
		return
	}

	link := _link[0]

	// 需要展示的数据
	show := make(map[string]interface{})

	// 所有类别/导航
	Categories := services.AllCategoryData()

	MovieDetail := services.MovieDetail(link)


	show["categories"] = Categories
	show["MovieDetail"] = MovieDetail

	show["nav_link"] = "/"
	show["film_title"] = MovieDetail["info"].(map[string]string)["name"]

	//log.Println(show)
	tmpl.GoTpl.ExecuteTemplate(w,"detail",show)
}

func Play(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	show := make(map[string]interface{})

	PlayUrl := r.URL.Query().Get("play_url")

	PlayType := "kuyun"
	if strings.Contains(PlayUrl, ".mp4"){
		PlayType = "mp4"
	} else if strings.Contains(PlayUrl, ".m3u8"){
		PlayType = "m3u8"
	}

	show["playUrl"] = PlayUrl
	show["playType"] = PlayType

	Categories := services.AllCategoryData()
	link := r.URL.Query().Get("link")
	episode := r.URL.Query().Get("episode")
	MovieDetail := services.MovieDetail(link)
	show["MovieDetail"] = MovieDetail
	show["categories"] = Categories
	show["film_title"] = MovieDetail["info"].(map[string]string)["name"]
	show["episode"] = episode

	tmpl.GoTpl.ExecuteTemplate(w,"play",show)
}

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	show := make(map[string]interface{})

	q := ""
	_q := r.URL.Query()["q"]
	if len(_q) > 0 {
		q = _q[0]
	}

	var MovieLists []services.MovieListStruct

	if strings.TrimSpace(q) != "" {
		MovieLists = services.SearchMovies(q)
	}

	// 所有类别/导航
	Categories := services.AllCategoryData()

	NewFilmKey := "detail_links:id:1"
	NewTVKey := "detail_links:id:2"

	NewFilm := services.MovieListsRange(NewFilmKey, 0, 49)
	NewTV := services.MovieListsRange(NewTVKey, 0, 49)
	recommend := services.MoviesRecommend()
	show["recommend"] = recommend
	show["movieLists"] = MovieLists
	show["categories"] = Categories
	show["new_film"] = NewFilm
	show["new_tv"] = NewTV
	show["nav_link"] = "/"
	show["film_title"] = ""

	buffer := new(bytes.Buffer)
	heroTpl.Search(show, buffer)
	w.Write(buffer.Bytes())
}

func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// 需要展示的数据
	show := make(map[string]interface{})

	// 所有类别/导航
	Categories := services.AllCategoryData()

	show["categories"] = Categories

	show["nav_link"] = "/about"
	show["film_title"] = ""

	buffer := new(bytes.Buffer)
	heroTpl.About(show, buffer)
	w.Write(buffer.Bytes())
}
