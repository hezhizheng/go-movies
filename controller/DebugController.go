package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_movies/utils/spider/tian_kong"
	"net/http"
)

func Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	go tian_kong.DelAllListCacheKey()

	fmt.Fprint(w, "DEBUG")
}
