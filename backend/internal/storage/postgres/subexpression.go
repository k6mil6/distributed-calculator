package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	"github.com/samber/lo"
)

type SubexpressionStorage struct {
	db *sqlx.DB
}

func NewSubExpressionStorage(db *sqlx.DB) *SubexpressionStorage {
	return &SubexpressionStorage{
		db: db,
	}
}

func (s *SubexpressionStorage) Save(context context.Context, subExpression model.Subexpression) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	row := conn.QueryRowContext(
		context,
		`INSERT INTO subexpressions (expression_id, subexpression, timeout) VALUES ($1, $2, $3)`,
		subExpression.ExpressionId,
		subExpression.Subexpression,
		subExpression.Timeout,
	)

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

func (s *SubexpressionStorage) NonTakenSubexpressions(context context.Context) ([]model.Subexpression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var subexpressions []dbSubexpression

	if err := conn.SelectContext(context,
		&subexpressions,
		`SELECT 
    	 id,
    	 expression_id,
    	 subexpression,
    	 timeout
    	 FROM subexpressions 
         WHERE is_taken = false`,
	); err != nil {
		return nil, err
	}

	return lo.Map(subexpressions, func(subexpression dbSubexpression, _ int) model.Subexpression {
		return model.Subexpression(subexpression)
	}), nil
}

func (s *SubexpressionStorage) TakeSubexpression(context context.Context, id, workerId int) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		context,
		`UPDATE subexpressions SET is_taken = true, worker_id = $1 WHERE id = $2`,
		workerId,
		id,
	); err != nil {
		return err
	}

	return nil
}

func (s *SubexpressionStorage) SubexpressionIsDone(context context.Context, id int, result float64) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		context,
		`UPDATE subexpressions SET is_done = true, result = $1 WHERE id = $2`,
		result,
		id,
	); err != nil {
		return err
	}

	return nil
}

type dbSubexpression struct {
	ID            int       `db:"id"`
	ExpressionId  uuid.UUID `db:"expression_id"`
	WorkerId      int       `db:"worker_id"`
	Subexpression string    `db:"subexpression"`
	IsTaken       bool      `db:"is_taken"`
	IsDone        bool      `db:"is_done"`
	Timeout       int64     `db:"timeout"`
	Result        float64   `db:"result"`
}
