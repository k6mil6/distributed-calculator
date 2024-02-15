package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	"github.com/k6mil6/distributed-calculator/backend/internal/storage"
	"github.com/k6mil6/distributed-calculator/backend/internal/timeout"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"time"
)

type ExpressionStorage struct {
	db *sqlx.DB
}

func NewExpressionStorage(db *sqlx.DB) *ExpressionStorage {
	return &ExpressionStorage{
		db: db,
	}
}

func (s *ExpressionStorage) Save(context context.Context, expression model.Expression) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	timeouts, err := json.Marshal(expression.Timeouts)
	if err != nil {
		return err
	}

	if _, err := conn.ExecContext(
		context,
		`INSERT INTO expressions (id, expression, timeouts) VALUES ($1, $2, $3)`,
		expression.ID,
		expression.Expression,
		timeouts,
	); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // PostgreSQL error code for unique_violation
				return storage.ErrExpressionInProgress
			}
			return err
		}
		return err
	}

	return nil
}

func (s *ExpressionStorage) Get(context context.Context, id uuid.UUID) (model.Expression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return model.Expression{}, err
	}
	defer conn.Close()

	var expression dbExpression

	if err := conn.GetContext(context, &expression, `SELECT * FROM expressions WHERE id = $1`, id); err != nil {
		return model.Expression{}, err
	}

	return model.Expression(expression), nil
}

func (s *ExpressionStorage) NonTakenExpressions(context context.Context) ([]model.Expression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var expressions []dbExpression

	if err := conn.SelectContext(context, &expressions, `SELECT id, expression, created_at, timeouts, is_taken, is_done FROM expressions WHERE is_taken = false ORDER BY created_at`); err != nil {
		return nil, err
	}

	return lo.Map(expressions, func(expression dbExpression, _ int) model.Expression {
		return model.Expression(expression)
	}), nil
}

func (s *ExpressionStorage) TakeExpression(context context.Context, id uuid.UUID) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		context,
		`UPDATE expressions SET is_taken = true WHERE id = $1`,
		id,
	)

	return err
}

func (s *ExpressionStorage) AllExpressions(context context.Context) ([]model.Expression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var expressions []dbExpression

	if err := conn.SelectContext(context, &expressions, `SELECT id, expression, created_at, is_taken, is_done, result FROM expressions ORDER BY created_at`); err != nil {
		return nil, err
	}

	return lo.Map(expressions, func(expression dbExpression, _ int) model.Expression {
		return model.Expression(expression)
	}), nil
}

func (s *ExpressionStorage) UpdateResult(context context.Context, id uuid.UUID, result float64) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		context,
		`UPDATE expressions SET is_done = true, result = $1 WHERE id = $2`,
		result,
		id,
	)

	return err
}

type dbExpression struct {
	ID         uuid.UUID       `db:"id"`
	Expression string          `db:"expression"`
	CreatedAt  time.Time       `db:"created_at"`
	Timeouts   timeout.Timeout `db:"timeouts"`
	IsTaken    bool            `db:"is_taken"`
	IsDone     bool            `db:"is_done"`
	Result     float64         `db:"result"`
}
