package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.wrote = true
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.wrote = true
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	} else {
		fmt.Println("http.Flusher unsupported for ResponseWriter")
	}
}

type rootHandler struct {
	mux *mux.Router
	s   *server
}
