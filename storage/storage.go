package storage

import (
	"context"
	"github.com/zhayt/kmf-tt/config"
	"github.com/zhayt/kmf-tt/model"
	"github.com/zhayt/kmf-tt/storage/mssql"
	"go.uber.org/zap"
	"time"
)

type ICurrencyStorage interface {
	SaveCurrency(ctx context.Context, currency model.Currency) error
	GetCurrencyByDate(ctx context.Context, dateTime time.Time) ([]model.Currency, error)
	GetCurrencyByDateCode(ctx context.Context, dateTime time.Time, code string) ([]model.Currency, error)
}

type Storage struct {
	ICurrencyStorage
}

func NewStorage(cfg *config.Config, l *zap.Logger) (*Storage, error) {
	db, err := mssql.Dial(cfg)
	if err != nil {
		return nil, err
	}

	return &Storage{mssql.NewStorage(db, l)}, nil
}
