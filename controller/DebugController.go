package controller

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

func Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//x := services.MovieListsRange("detail_links:id:2", 0, 15)

	log.Println(strconv.FormatInt(20-14, 10))
}
