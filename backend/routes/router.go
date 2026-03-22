package routes

import (
	"net/http"
)

func Routes(mux *http.ServeMux, handler *Handler) {
	Routes := map[string]http.HandlerFunc{
		"/api/login": handler.Cntrlrs.Login,
	}

	for path, h := range Routes {
		mux.HandleFunc(path, h)
	}
}
