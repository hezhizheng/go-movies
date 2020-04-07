package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_movies/utils/spider"
	"net/http"
)

func Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	go spider.DelAllListCacheKey()

	fmt.Fprint(w, "DEBUG")
}
