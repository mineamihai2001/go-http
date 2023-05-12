package core

import (
	"fmt"
	"regexp"
	"strings"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

type Route struct {
	path   string
	method HttpMethod
}

type RouteHandler func(req *Request, res *Response) // w http.ResponseWriter,

type Router struct {
	routes map[Route]RouteHandler
}

func NewRouter() Router {
	return Router{routes: make(map[Route]RouteHandler)}
}

func (r *Router) register(route Route, callback RouteHandler) {
	r.routes[route] = callback
}

func (r *Router) get(route Route) (RouteHandler, error) {
	var re = regexp.MustCompile(`:([^\/]+)`)
	routeArgs := Filter(strings.Split(route.path, "/"), func(item string) bool {
		return item != ""
	})
	for r, c := range r.routes {
		stripped := re.FindAllString(r.path, -1)
		rArgs := Filter(strings.Split(r.path, "/"), func(item string) bool {
			return item != ""
		})

		if r.method == route.method && (r.path == route.path || (len(stripped) > 0 && len(rArgs) == len(routeArgs))) {
			return c, nil
		}
	}
	return nil, &HttpError{
		StatusCode: 404,
		Err:        fmt.Errorf("not found: %v", route),
	}
}
