package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) InitRoute() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("api/v1/currency/save/{date}", s.handler.SaveCurrency).Methods("GET")
	r.HandleFunc("api/v1/currency/{date}/{code}", s.handler.ShowCurrency).Methods("GET")

	return r
}
