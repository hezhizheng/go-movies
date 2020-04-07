package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_movies/utils/spider"
	"net/http"
)

func GoSpider(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//defer ants.Release()
	//go utils.StartSpider()

	//go utils.MoviesInfo("/?m=vod-detail-id-41731.html")

	go spider.StartApi()

	fmt.Fprint(w, "Spider ing....")
}
