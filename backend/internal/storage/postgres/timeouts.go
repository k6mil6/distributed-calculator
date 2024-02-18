package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	"github.com/k6mil6/distributed-calculator/backend/internal/timeout"
)

type TimeoutsStorage struct {
	db *sqlx.DB
}

func NewTimeoutsStorage(db *sqlx.DB) *TimeoutsStorage {
	return &TimeoutsStorage{
		db: db,
	}
}

func (s *TimeoutsStorage) Save(context context.Context, timeouts model.Timeouts) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		context,
		`INSERT INTO timeouts (timeouts_values) VALUES ($1)`,
		timeouts.Timeouts,
	); err != nil {
		return err
	}

	return nil
}

func (s *TimeoutsStorage) GetActualTimeouts(context context.Context) (model.Timeouts, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return model.Timeouts{}, err
	}
	defer conn.Close()

	var timeouts dbTimeouts

	if err := conn.GetContext(
		context,
		&timeouts,
		`SELECT timeouts_values FROM timeouts ORDER BY id DESC LIMIT 1`,
	); err != nil {
		return model.Timeouts{}, err
	}
	return model.Timeouts(timeouts), nil
}

type dbTimeouts struct {
	ID       int             `db:"id"`
	Timeouts timeout.Timeout `db:"timeouts_values"`
}
