package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	"github.com/samber/lo"
	"time"
)

type SubExpressionStorage struct {
	db *sqlx.DB
}

func NewSubExpressionStorage(db *sqlx.DB) *SubExpressionStorage {
	return &SubExpressionStorage{
		db: db,
	}
}

func (s *SubExpressionStorage) Save(context context.Context, subExpression model.Subexpression) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	row := conn.QueryRowContext(
		context,
		`INSERT INTO subexpressions (expression_id, subexpression, is_taken) VALUES ($1, $2, $3)`,
		subExpression.ExpressionId,
		subExpression.Subexpression,
		subExpression.IsTaken,
	)

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

func (s *SubExpressionStorage) NonTakenSubexpressions(context context.Context) ([]model.Subexpression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var subexpressions []dbSubexpression
	if err := conn.SelectContext(context, &subexpressions, `SELECT * FROM subexpressions WHERE is_taken = false`); err != nil {
		return nil, err
	}

	return lo.Map(subexpressions, func(subexpression dbSubexpression, _ int) model.Subexpression {
		return model.Subexpression(subexpression)
	}), nil
}

func (s *ExpressionPostgresStorage) TakeSubexpression(context context.Context, id int) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	row := conn.QueryRowContext(
		context,
		`UPDATE subexpressions SET is_taken = true WHERE id = $1`,
		id,
	)

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

func (s *ExpressionPostgresStorage) SubexpressionIsDone(context context.Context, id int) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	row := conn.QueryRowContext(
		context,
		`UPDATE subexpressions SET is_done = true WHERE id = $1`,
		id,
	)

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

type dbSubexpression struct {
	ID            int           `db:"id"`
	ExpressionId  uuid.UUID     `db:"expression_id"`
	WorkerId      int           `db:"worker_id"`
	Subexpression string        `db:"subexpression"`
	IsTaken       bool          `db:"is_taken"`
	IsDone        bool          `db:"is_done"`
	Timeout       time.Duration `db:"timeout"`
	Result        float64       `db:"result"`
}
