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
	GetCurrency(ctx context.Context, date string, code string) ([]model.Currency, error)
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
		s.l.Error("Parse date error", zap.Error(err))
		return fmt.Errorf("failde to parse date: %w", ErrUserStupid)
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

		go func() {
			if err := s.currency.SaveCurrency(ctx, currency); err != nil {
				s.l.Error("SaveCurrency error", zap.Error(err))
			}
		}()
	}

	return nil
}

func (s *CurrencyService) GetCurrency(ctx context.Context, date string, code string) ([]model.Currency, error) {
	// validate data
	dateTime, err := time.Parse("02.01.2006", date)
	if err != nil {
		s.l.Error("Parse date error", zap.Error(err))
		return []model.Currency{}, fmt.Errorf("failde to parse date: %w", ErrUserStupid)
	}

	if code == "" {
		return s.currency.GetCurrencyByDate(ctx, dateTime)
	}

	return s.currency.GetCurrencyByDateCode(ctx, dateTime, code)
}

func makeURL(api, date string) string {
	builder := strings.Builder{}

	builder.WriteString(api)
	builder.WriteString(date)

	return builder.String()
}
