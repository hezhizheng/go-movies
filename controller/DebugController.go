package controller

import (
	"fmt"
	"go_movies/utils/spider/haohua"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	go haohua.DelAllListCacheKey()

	fmt.Fprint(w, "DEBUG")
}
