package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/zhayt/kmf-tt/model"
	"github.com/zhayt/kmf-tt/storage"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type ICurrencyService interface {
	SaveCurrency(ctx context.Context, date string) error
	GetCurrencyByDate(ctx context.Context, date string) ([]model.Currency, error)
	GetCurrencyByDateCode(ctx context.Context, date string, code string) ([]model.Currency, error)
}

type CurrencyService struct {
	currency    *storage.Storage
	l           *zap.Logger
	externalAPI string
}

func NewCurrencyService(currency *storage.Storage, l *zap.Logger) *CurrencyService {
	api := "https://nationalbank.kz/rss/get_rates.cfm?fdate="
	return &CurrencyService{currency: currency, l: l, externalAPI: api}
}

func (s *CurrencyService) SaveCurrency(ctx context.Context, date string) error {
	// validate data
	dateTime, err := time.Parse("02.01.2006", date)
	if err != nil {
		return fmt.Errorf("%w, failde to parse date: %w", ErrUserStupid, err)
	}

	// work with external api
	url := makeURL(s.externalAPI, date)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to retrieve data from external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get data from external API, received status code %d", resp.StatusCode)
	}

	// parse xml
	var rates Rates

	if err = xml.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return fmt.Errorf("failed to decode xml: %w", err)
	}

	for _, item := range rates.Items {
		currency := model.Currency{
			Title: item.FullName,
			Code:  item.Title,
			Value: item.Description,
			ADate: dateTime,
		}

		go s.currency.SaveCurrency(ctx, currency)
	}

	return nil
}

func makeURL(api, date string) string {
	builder := strings.Builder{}

	builder.WriteString(api)
	builder.WriteString(date)

	return builder.String()
}
