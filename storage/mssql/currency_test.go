package mssql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/kmf-tt/model"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func TestCurrencyStorage_SaveCurrency1(t *testing.T) {
	type fields struct {
		db *sqlx.DB
		l  *zap.Logger
	}
	type args struct {
		ctx      context.Context
		currency model.Currency
	}

	dbContainer, db, err := SetupTestDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer dbContainer.Terminate(context.Background())

	l := zap.NewExample()

	fiel := fields{db: db, l: l}

	dateTime, _ := time.Parse("02.01.2006", "24.01.2002")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fiel, args: args{context.Background(), model.Currency{Title: "фыафы", Code: "USE", Value: 145.21, ADate: dateTime}}},
		{"failed lond code", fiel, args{context.Background(), model.Currency{Title: "ASD", Code: "USEA", Value: 145.21, ADate: dateTime}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CurrencyStorage{
				db: tt.fields.db,
				l:  tt.fields.l,
			}
			if err := r.SaveCurrency(tt.args.ctx, tt.args.currency); (err != nil) != tt.wantErr {
				t.Errorf("SaveCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
