package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/model"
)

type DB struct {
	dbPool *pgxpool.Pool
}

func (p *DB) Ping(ctx context.Context) error {
	return p.dbPool.Ping(ctx)
}

func (p *DB) Init(ctx context.Context) error {
	var err error
	p.dbPool, err = pgxpool.New(ctx, flags.AddrDB)
	if err != nil {
		return err
	}
	return p.NewPostgresStorage(ctx)
}

func (p *DB) NewPostgresStorage(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS metrics (
	    id VARCHAR(255) PRIMARY KEY,
	    type VARCHAR(255) NOT NULL DEFAULT '',
	    delta INTEGER NOT NULL DEFAULT 0,
	    value DOUBLE PRECISION NOT NULL DEFAULT 0,
	    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := p.dbPool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (p *DB) SaveMetric(ctx context.Context, m model.Metrics) error {
	switch m.MType {
	case "gauge":
		_, err := p.dbPool.Exec(ctx, `INSERT INTO public.metrics (id, type, value)
												VALUES ($1,$2,$3)
												ON CONFLICT (id)
												DO UPDATE SET type = excluded.type, value = excluded.value;`,
			m.ID, m.MType, *m.Value)
		if err != nil {
			return err
		}
	case "counter":
		_, err := p.dbPool.Exec(ctx, `INSERT INTO public.metrics (id, type, delta)
												VALUES ($1,$2,$3)
												ON CONFLICT (id)
												DO UPDATE SET delta = (metrics.delta + excluded.delta);`,
			m.ID, m.MType, *m.Delta)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *DB) GetAllMetrics(ctx context.Context) (map[string]model.Metrics, error) {
	row, err := p.dbPool.Query(ctx, `SELECT id, type, delta, value FROM public.metrics`)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	res := make(map[string]model.Metrics)

	for row.Next() {
		var metric model.Metrics

		if err = row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
			return nil, err
		}
		res[metric.ID] = metric
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *DB) GetByID(ctx context.Context, id string) (model.Metrics, error) {
	row := p.dbPool.QueryRow(ctx, `SELECT id, type, delta, value FROM public.metrics WHERE id = $1`, id)
	var m model.Metrics
	err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
	if err != nil {
		return model.Metrics{}, err
	}
	return m, nil
}

func (p *DB) Close() {
	p.dbPool.Close()
}
