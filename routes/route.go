package routes

import (
	"github.com/julienschmidt/httprouter"
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
	routes := Routes{
		Route{"Index", "GET", "/", controller.Index},
		Route{"Movie", "GET", "/movie", controller.Movie},
		Route{"Search", "GET", "/search", controller.Search},
		Route{"Play", "GET", "/play", controller.Play},
		Route{"About", "GET", "/about", controller.About},
		Route{"Debug", "GET", "/debug", controller.Debug},
		Route{"GoSpider", "GET", "/go-spider", controller.GoSpider},
	}
	return routes
}
