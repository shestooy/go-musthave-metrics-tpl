package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/model"
	"log"
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
	    value DOUBLE PRECISION NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_id ON metrics (id)`

	_, err := p.dbPool.Exec(ctx, query)
	return err
}

func (p *DB) SaveMetric(ctx context.Context, m model.Metrics) (model.Metrics, error) {
	switch m.MType {
	case gauge:
		err := p.dbPool.QueryRow(ctx, `INSERT INTO metrics (id, type, value)
												VALUES ($1,$2,$3)
												ON CONFLICT (id)
												DO UPDATE SET type = excluded.type, value = excluded.value
												RETURNING id, type, value;`,
			m.ID, m.MType, *m.Value).Scan(&m.ID, &m.MType, &m.Value)
		if err != nil {
			return m, err
		}
	case counter:
		err := p.dbPool.QueryRow(ctx, `INSERT INTO metrics (id, type, delta)
												VALUES ($1,$2,$3)
												ON CONFLICT (id)
												DO UPDATE SET delta = (metrics.delta + excluded.delta)
												RETURNING id, type, delta;`,
			m.ID, m.MType, *m.Delta).Scan(&m.ID, &m.MType, &m.Delta)
		if err != nil {
			return m, err
		}
	}

	return m, nil
}

func (p *DB) GetAllMetrics(ctx context.Context) (map[string]model.Metrics, error) {
	row, err := p.dbPool.Query(ctx, `SELECT id, type, delta, value FROM metrics`)
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
	row := p.dbPool.QueryRow(ctx, `SELECT id, type, delta, value FROM metrics WHERE id = $1`, id)
	var m model.Metrics

	err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
	if err != nil {
		return model.Metrics{}, err
	}
	switch m.MType {
	case gauge:
		m.Delta = nil
	case counter:
		m.Value = nil
	}

	return m, nil
}

func (p *DB) SaveMetrics(ctx context.Context, metrics []model.Metrics) ([]model.Metrics, error) {
	tx, err := p.dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	ans := make([]model.Metrics, 0)
	for _, metric := range metrics {
		var ansMetric model.Metrics
		switch metric.MType {
		case gauge:
			err = tx.QueryRow(ctx, `INSERT INTO metrics (id, type, value) VALUES ($1, $2, $3)
									ON CONFLICT (id)
									DO UPDATE SET value = excluded.value
									RETURNING value`, metric.ID, metric.MType, *metric.Value).Scan(
				&ansMetric.Value)
			if err != nil {
				log.Println(err.Error())
				return nil, err
			}
			ansMetric.MType = metric.MType
			ansMetric.ID = metric.ID
			ans = append(ans, ansMetric)
		case counter:
			err = tx.QueryRow(ctx, `INSERT INTO metrics (id, type, delta) VALUES ($1, $2, $3)
									ON CONFLICT (id)
									DO UPDATE SET delta = metrics.delta +excluded.delta
									RETURNING delta`, metric.ID, metric.MType, *metric.Delta).Scan(&ansMetric.Delta)
			if err != nil {
				log.Println(err.Error())
				return nil, err
			}
			ansMetric.MType = metric.MType
			ansMetric.ID = metric.ID
			ans = append(ans, ansMetric)
		}
	}
	err = tx.Commit(ctx)
	return ans, err
}

func (p *DB) Close() {
	p.dbPool.Close()
}
