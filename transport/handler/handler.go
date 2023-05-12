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

	if err := h.currency.SaveCurrency(context.TODO(), date); err != nil {

		if errors.Is(err, service.ErrUserStupid) {
			h.respondWithError(w, http.StatusBadRequest, "The entered date is not valid.")
			return
		}

		h.respondWithError(w, http.StatusInternalServerError, "Couldn't save currency")
		return
	}

	h.respondWithSuccess(w, http.StatusOK)
}

func (h *Handler) ShowCurrency(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.l.Error(message)

	response := model.ErrorResponse{
		Message: message,
	}

	json, err := formatToJSON(response)
	if err != nil {
		h.l.Error("formatToJSON error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something wrong try later!"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}

func (h *Handler) respondWithSuccess(w http.ResponseWriter, code int) {
	response := model.SuccessResponse{
		Success: true,
	}

	json, err := formatToJSON(response)
	if err != nil {
		h.l.Error("formatToJSON error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something wrong try later!"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}

func formatToJSON(input interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	switch v := input.(type) {
	case model.SuccessResponse:
		return json.Marshal(input.(model.SuccessResponse))
	case model.ErrorResponse:
		return json.Marshal(input.(model.ErrorResponse))
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}
