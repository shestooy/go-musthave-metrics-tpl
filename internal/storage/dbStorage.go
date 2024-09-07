package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
)

var DBPool *pgxpool.Pool

func InitDB(ctx context.Context) error {
	var err error
	DBPool, err = pgxpool.New(ctx, flags.AddrDB)
	if err != nil {
		return err
	}
	return nil
}
