package http

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

func (s *Server) InitRoute() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/currency/save/{date}", s.handler.SaveCurrency).Methods("GET")
	r.HandleFunc("/api/v1/currency/{date}", s.handler.ShowCurrency).Methods("GET")
	r.HandleFunc("/api/v1/currency/{date}/{code}", s.handler.ShowCurrency).Methods("GET")

	return r
}
