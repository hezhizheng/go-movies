package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_movies/utils/spider"
	"net/http"
)

func GoSpider(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	spider.Create().Start()

	fmt.Fprint(w, "Spider ing....")
}
