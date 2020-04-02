package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"go_movies/controller"
)

/*
Define all the routes here.
A new Route entry passed to the routes slice will be automatically
translated to a handler with the NewRouter() function
*/
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func AllRoutes() Routes {
	spiderPath := viper.GetString(`app.spider_path`)
	spiderPathName := viper.GetString(`app.spider_path_name`)

	debugPath := viper.GetString(`app.debug_path`)
	debugPathName := viper.GetString(`app.debug_path_name`)
	routes := Routes{
		Route{"Index", "GET", "/", controller.Index},
		Route{"Movie", "GET", "/movie", controller.Movie},
		Route{"Search", "GET", "/search", controller.Search},
		Route{"Play", "GET", "/play", controller.Play},
		Route{"About", "GET", "/about", controller.About},
		Route{debugPathName, "GET", debugPath, controller.Debug},
		Route{spiderPathName, "GET", spiderPath, controller.GoSpider},
	}
	return routes
}
