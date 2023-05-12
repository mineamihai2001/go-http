package core

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Request struct {
	HttpRequest *http.Request
}

// parses the body into map[any]any
func (r *Request) Body() any {
	body, _ := ioutil.ReadAll(r.HttpRequest.Body)
	var temp any
	Parse(string(body), &temp)
	return temp
}

// parses the body into a <T> object
func Body[T any](req *Request, out *T) {
	body, _ := ioutil.ReadAll(req.HttpRequest.Body)
	Parse(string(body), out)
}

// parses the query params into a map[string][]string
func (r *Request) Query() map[string][]string {
	values := r.HttpRequest.URL.Query()
	result := make(map[string][]string)
	for k, v := range values {
		result[k] = v
	}
	return result
}

func (r *Request) Params() []string {
	path := r.HttpRequest.URL.Path
	params := strings.Split(path, "/")
	var result []string
	for _, param := range params {
		if param != "" {
			result = append(result, param)
		}
	}

	return result
}
