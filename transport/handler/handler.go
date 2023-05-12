package handler

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhayt/kmf-tt/service"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	currency *service.CurrencyService
	l        *zap.Logger
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewHandler(currency *service.CurrencyService, l *zap.Logger) *Handler {
	return &Handler{currency: currency, l: l}
}

func (h *Handler) SaveCurrency(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]

	h.l.Info("Get url param", zap.String("date", date))

	if err := h.currency.SaveCurrency(context.TODO(), date); err != nil {

		h.l.Error("SaveCurrency error", zap.Error(err))

		if errors.Is(err, service.ErrUserStupid) {
			h.respondWithError(w, http.StatusBadRequest, "The entered date is not valid.")
			return
		}

		h.respondWithError(w, http.StatusInternalServerError, "Couldn't save currency")
		return
	}

	h.l.Info("Currency saved", zap.String("date", date))

	h.respondWithSuccess(w, http.StatusOK)
}

func (h *Handler) ShowCurrency(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]
	code := mux.Vars(r)["code"]

	h.l.Info("Get url param", zap.String("date", date), zap.String("code", code))

	// rest of the handler code
	currencies, err := h.currency.GetCurrency(context.TODO(), date, code)
	if err != nil {
		h.l.Error("GetCurrency error", zap.Error(err))
		if errors.Is(err, service.ErrUserStupid) {
			h.respondWithError(w, http.StatusBadRequest, "The entered date is not valid.")
			return
		}

		h.respondWithError(w, http.StatusInternalServerError, "Couldn't get currencies")
		return
	}

	if len(currencies) == 0 {
		h.l.Error("Zero currencies error", zap.Error(err))
		h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	h.l.Info("Currencies founded", zap.Int("amount", len(currencies)))

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	json.NewEncoder(w).Encode(currencies)
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := ErrorResponse{
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) respondWithSuccess(w http.ResponseWriter, code int) {
	response := SuccessResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	json.NewEncoder(w).Encode(response)
}
