package main

import (
	"github.com/julienschmidt/httprouter"
	"go_movies/routes"
	"log"
	"net/http"
)

// Reads from the routes slice to translate the values to httprouter.Handle
// 遍历路由
func TraversingRouter() *httprouter.Router {

	AllRoutes := routes.AllRoutes()

	router := httprouter.New()
	for _, route := range AllRoutes {
		var handle httprouter.Handle

		handle = route.HandlerFunc

		log.Println(route.Path)
		router.Handle(route.Method, route.Path, handle)
	}

	// 配置静态文件访问
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	return router
}

func main() {
	// 注册所有路由
	router := TraversingRouter()

	addr := ":8899"
	log.Println("监听端口", "http://127.0.0.1"+addr)
	log.Println(http.ListenAndServe(addr, router))
}
