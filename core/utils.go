package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Err        error
}

func (r *HttpError) Error() string {
	return fmt.Sprintf("status: %v: err: %v", r.StatusCode, r.Err)
}

func Check(err error, args ...string) {
	message := ""
	for _, arg := range args {
		message += arg
	}
	if err != nil {
		panic(fmt.Sprintf("%v: %v", message, err))
	}
}

func HttpCheck(w http.ResponseWriter, err error, args ...string) {
	message := ""
	for _, arg := range args {
		message += arg
	}
	if err != nil {
		w.WriteHeader(404)
		io.WriteString(w, fmt.Sprintf("%v: %v", message, err))
	}
}

// Deserializes a json into the given out object
func Parse[T any, V []byte | string](data V, out *T) {
	json.Unmarshal([]byte(data), out)
}

// Deserializes a json into the given out object
func ToJSON[T any, V []byte | string](data V, out *T) *T {
	res := out
	json.Unmarshal([]byte(data), res)
	return res
}

// Serializes a struct into a json
func Stringify[T any](object T) string {
	bytes, err := json.Marshal(object)
	Check(err)

	return string(bytes)
}

// Copies the values of a list in another
// Returns a cloned list with the values from the initial
func Clone[T any](list []T) []T {
	cpy := make([]T, len(list))
	copy(cpy, list)
	return cpy
}

func Filter[T any](source []T, filter func(x T) bool) []T {
	var result []T
	for _, item := range source {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any](source []T, _map func(x T) T) []T {
	var result []T
	for _, item := range source {
		result = append(result, _map(item))
	}
	return result
}
