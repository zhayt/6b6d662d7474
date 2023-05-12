package mssql

import (
	"context"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/kmf-tt/config"
	"time"
)

const _defaultTimeout = 10 * time.Second

func Dial(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("couldn't get pool connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), _defaultTimeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("couldn't connect db: %w", err)
	}

	return db, nil
}
