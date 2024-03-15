package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/AvyChanna/nginx-token-authz/lib/rbac"
)

func ping(autherPtr *atomic.Pointer[rbac.Auther]) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		val, err := json.Marshal(autherPtr.Load())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(val)
	}
}

func NewMux(autherPtr *atomic.Pointer[rbac.Auther]) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", ping(autherPtr))

	return mux
}
