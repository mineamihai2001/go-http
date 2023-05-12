package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	httpServer *http.Server
	port       int
	addr       string
	router     Router
}

func NewServer(port int) Server {
	address := fmt.Sprintf(":%v", port)
	server := Server{
		router: NewRouter(),
		addr:   address,
		port:   port,
	}
	httpServer := &http.Server{
		Addr:         address,
		Handler:      http.HandlerFunc(server.serverHandler),
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
	server.httpServer = httpServer
	return server
}

func (s *Server) Run(callback func()) {
	callback()
	log.Fatal(s.httpServer.ListenAndServe())
}

func (s *Server) serverHandler(writer http.ResponseWriter, req *http.Request) {
	currentRoute := Route{
		path:   req.URL.Path,
		method: HttpMethod(strings.ToUpper(req.Method)),
	}
	fmt.Printf(">>> %v %v\n", currentRoute.method, currentRoute.path)

	// create the response object
	res := &Response{
		Headers: make(map[string]string),
	}

	// create the request object
	request := &Request{
		HttpRequest: req,
	}

	// validate the route
	callback, err := s.router.get(currentRoute)
	if err != nil {
		// throw 404 error
		HttpCheck(writer, err)
		return
	}
	
	// execute the registered callback
	callback(request, res)

	// write the header
	res.WriteHeaders(writer)
	// write the response
	io.WriteString(writer, res.body)
}

/**************************************
* 			HTTP METHODS
***************************************/

func (s *Server) Get(path string, callback func(r *Request, res *Response)) {
	route := Route{
		path:   path,
		method: GET,
	}
	s.router.register(route, callback)
}

func (s *Server) Post(path string, callback func(r *Request, res *Response)) {
	route := Route{
		path:   path,
		method: POST,
	}
	s.router.register(route, callback)
}

func (s *Server) Put(path string, callback func(r *Request, res *Response)) {
	route := Route{
		path:   path,
		method: PUT,
	}
	s.router.register(route, callback)
}

func (s *Server) Delete(path string, callback func(r *Request, res *Response)) {
	route := Route{
		path:   path,
		method: DELETE,
	}
	s.router.register(route, callback)
}
