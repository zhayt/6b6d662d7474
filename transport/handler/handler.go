package handler

import (
	"github.com/gorilla/mux"
	"github.com/zhayt/kmf-tt/service"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	currency *service.CurrencyService
	l        *zap.Logger
}

func NewHandler(currency *service.CurrencyService, l *zap.Logger) *Handler {
	return &Handler{currency: currency, l: l}
}

func (h *Handler) SaveCurrency(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]

	h.l.Info("Get url param", zap.String("date", date))

}

func (h *Handler) ShowCurrency(w http.ResponseWriter, r *http.Request) {

}
