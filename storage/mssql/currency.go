package mssql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/kmf-tt/model"
	"go.uber.org/zap"
	"time"
)

type CurrencyStorage struct {
	db *sqlx.DB
	l  *zap.Logger
}

func NewStorage(db *sqlx.DB, l *zap.Logger) *CurrencyStorage {
	return &CurrencyStorage{db: db, l: l}
}

func (r *CurrencyStorage) SaveCurrency(ctx context.Context, currency model.Currency) error {
	qr := `INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, qr, currency.Title, currency.Code, currency.Value, currency.ADate)
	if err != nil {
		return fmt.Errorf("couldn't save currency: %w", err)
	}

	return nil
}

func (r *CurrencyStorage) GetCurrencyByDate(ctx context.Context, date time.Time) ([]model.Currency, error) {
	qr := `SELECT * FROM R_CURRENCY WHERE A_DATE = ?`

	var currencies []model.Currency

	if err := r.db.SelectContext(ctx, &currencies, qr, date); err != nil {
		return []model.Currency{}, err
	}

	return currencies, nil
}

func (r *CurrencyStorage) GetCurrencyByDateCode(ctx context.Context, date time.Time, code string) ([]model.Currency, error) {
	qr := `SELECT * FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?`

	var currencies []model.Currency

	if err := r.db.SelectContext(ctx, &currencies, qr, date, code); err != nil {
		return []model.Currency{}, err
	}

	return currencies, nil
}
