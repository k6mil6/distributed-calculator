package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	"github.com/samber/lo"
	"time"
)

type ExpressionPostgresStorage struct {
	db *sqlx.DB
}

func NewExpressionStorage(db *sqlx.DB) *ExpressionPostgresStorage {
	return &ExpressionPostgresStorage{
		db: db,
	}
}

func (s *ExpressionPostgresStorage) Save(context context.Context, expression model.Expression) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	row := conn.QueryRowContext(
		context,
		`INSERT INTO expressions (id, expression, is_taken) VALUES ($1, $2, $3)`,
		expression.ID,
		expression.Expression,
	)

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

func (s *ExpressionPostgresStorage) Get(context context.Context, id uuid.UUID) (model.Expression, error) {
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

func (s *ExpressionPostgresStorage) NonTakenExpressions(context context.Context) ([]model.Expression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var expressions []dbExpression

	if err := conn.SelectContext(context, &expressions, `SELECT * FROM expressions WHERE is_taken = false`); err != nil {
		return nil, err
	}

	return lo.Map(expressions, func(expression dbExpression, _ int) model.Expression {
		return model.Expression(expression)
	}), nil
}

func (s *ExpressionPostgresStorage) TakeExpression(context context.Context, id uuid.UUID) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	row := conn.QueryRowContext(
		context,
		`UPDATE expressions SET is_taken = true WHERE id = $1`,
		id,
	)

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

type dbExpression struct {
	ID         uuid.UUID `db:"id"`
	Expression string    `db:"expression"`
	CreatedAt  time.Time `db:"created_at"`
	IsTaken    bool      `db:"is_taken"`
	IsDone     bool      `db:"is_done"`
}
