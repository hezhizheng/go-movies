package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/panjf2000/ants/v2"
	"go_movies/utils"
	"net/http"
)

func GoSpider(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	defer ants.Release()
	go utils.StartSpider()

	//go utils.MoviesInfo("/?m=vod-detail-id-41731.html")

	fmt.Fprint(w, "Spider ing....")
}
