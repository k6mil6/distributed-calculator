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
		`INSERT INTO subexpressions (expression_id, subexpression, timeout, depends_on) VALUES ($1, $2, $3, $4)`,
		subExpression.ExpressionId,
		subExpression.Subexpression,
		subExpression.Timeout,
		subExpression.DependsOn,
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

	//if err := conn.SelectContext(context,
	//	&subexpressions,
	//	`SELECT
	//	 id,
	//	 expression_id,
	//	 subexpression,
	//	 timeout
	//	 FROM subexpressions
	//     WHERE is_taken = false`,
	//); err != nil {
	//	return nil, err
	//}

	if err := conn.SelectContext(context,
		&subexpressions,
		`WITH non_dependent_subexpressions AS (
    			SELECT *
    			FROM subexpressions
    			WHERE depends_on IS NULL AND is_taken = false
				)
				SELECT * FROM non_dependent_subexpressions
				UNION
-- Select all expressions that depend on evaluated expressions
				SELECT e.*
				FROM subexpressions e
				JOIN subexpressions d ON e.depends_on = d.id
				WHERE d.result IS NOT NULL
				EXCEPT
-- Exclude expressions that are already evaluated
				SELECT e.*
				FROM subexpressions e
				JOIN subexpressions d ON e.depends_on = d.id
				WHERE d.result IS NOT NULL;`,
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
		`UPDATE subexpressions SET result = $1 WHERE id = $2`,
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
	Timeout       int64     `db:"timeout"`
	DependsOn     []int     `db:"depends_on"`
	Result        float64   `db:"result"`
}
