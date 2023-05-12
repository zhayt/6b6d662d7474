package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/zhayt/kmf-tt/config"
	"github.com/zhayt/kmf-tt/storage"
	"go.uber.org/zap"
	"testing"
)

func TestCurrencyService_SaveCurrency(t *testing.T) {
	type fields struct {
		currency    *storage.Storage
		l           *zap.Logger
		externalAPI string
	}
	type args struct {
		ctx  context.Context
		date string
	}

	cfg := &config.Config{}
	cfg.Database.Driver = "mssql"
	cfg.Database.DSN = "sqlserver://sa:kursPswdsuper1@localhost:1433?database=TEST&connection+timeout=30"

	l := zap.NewExample()

	repo, err := storage.NewStorage(cfg, l)
	require.NoError(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success",
			fields: fields{currency: repo, l: l, externalAPI: "https://nationalbank.kz/rss/get_rates.cfm?fdate="},
			args:   args{context.TODO(), "15.04.2021"}},
		{name: "invalid date",
			fields: fields{currency: repo, l: l, externalAPI: "https://nationalbank.kz/rss/get_rates.cfm?fdate="},
			args:   args{context.TODO(), "15.44.2021"}, wantErr: true},
		{name: "invalid date",
			fields: fields{currency: repo, l: l, externalAPI: "https://nationalbank.kz/rss/get_rates.cfm?fdate="},
			args:   args{context.TODO(), "asd"}, wantErr: true},
		{name: "success",
			fields: fields{currency: repo, l: l, externalAPI: "https://nationalbank.kz/rss/get_rates.cfm?fdate="},
			args:   args{context.TODO(), "15.04.2025"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CurrencyService{
				currency:    tt.fields.currency,
				l:           tt.fields.l,
				externalAPI: tt.fields.externalAPI,
			}
			if err := s.SaveCurrency(tt.args.ctx, tt.args.date); (err != nil) != tt.wantErr {
				t.Errorf("SaveCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(err)
		})
	}
}
