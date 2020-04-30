package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_movies/utils"
	"net/http"
)

func GoSpider(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	go utils.StartSpider()

	//go utils.MoviesInfo("/?m=vod-detail-id-41731.html")

	fmt.Fprint(w, "Spider ing....")
}
