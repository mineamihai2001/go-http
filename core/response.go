package core

import "net/http"

type Response struct {
	Headers map[string]string
	body    string
	status  int
}

func (r *Response) WriteHeaders(w http.ResponseWriter) {
	for key, value := range r.Headers {
		w.Header().Set(key, value)
		code := r.status
		if code == 0 {
			code = http.StatusOK
		}
		w.WriteHeader(code)
	}
}

func (r *Response) Raw(object string) *Response {
	r.Headers["Content-Type"] = "text/plain"
	r.body = object
	return r
}

func (r *Response) Json(object any) *Response {
	r.Headers["Content-Type"] = "application/json"
	r.body = Stringify(object)
	return r
}

func (r *Response) Html(object any) *Response {
	r.Headers["Content-Type"] = "text/html"
	r.body = Stringify(object)
	return r
}

func (r *Response) Status(status int) *Response {
	r.status = status
	return r
}
