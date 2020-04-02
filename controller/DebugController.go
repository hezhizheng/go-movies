package controller

import (
	"github.com/julienschmidt/httprouter"
	"go_movies/utils"
	"net/http"
)

func Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	go utils.DelAllListCacheKey()
}
