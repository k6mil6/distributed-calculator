package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	"github.com/lib/pq"
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

	if _, err := conn.ExecContext(
		context,
		`INSERT INTO subexpressions (id ,expression_id, subexpression, depends_on, timeout) VALUES ($1, $2, $3, $4, $5)`,
		subExpression.ID,
		subExpression.ExpressionId,
		subExpression.Subexpression,
		pq.Array(subExpression.DependsOn),
		subExpression.Timeout,
	); err != nil {
		fmt.Println(err)
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
	    WHERE is_taken = false AND depends_on IS NULL`,
	); err != nil {
		return nil, err
	}

	//if err := conn.SelectContext(context,
	//	&subexpressions,
	//	`WITH RECURSIVE evaluated_subexpressions AS (
	//			SELECT id, subexpression, depends_on, result, 0 as level
	//			FROM subexpressions
	//			WHERE depends_on IS NULL AND result IS NULL AND is_taken = false -- Non-dependent and not evaluated
	//
	//			UNION ALL
	//
	//			SELECT s.id, s.subexpression, s.depends_on, s.result, ee.level + 1
	//			FROM subexpressions s
	//			INNER JOIN evaluated_subexpressions ee ON s.depends_on = ee.id
	//			WHERE s.result IS NULL AND s.is_taken = false AND ee.result IS NOT NULL  -- Dependent on evaluated and not evaluated itself
	//		)
	//		SELECT id, subexpression, depends_on, level FROM evaluated_subexpressions
	//		ORDER BY level;`,
	//); err != nil {
	//	return nil, err
	//}

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

func (s *SubexpressionStorage) LastSubexpression(context context.Context) (int, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int

	if err := conn.GetContext(context, &id, `SELECT id FROM subexpressions ORDER BY id DESC LIMIT 1`); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SubexpressionStorage) DoneSubexpressions(context context.Context) ([]model.Subexpression, error) {
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
		 timeout,
		 result
		 FROM subexpressions
	    WHERE result IS NOT NULL`,
	); err != nil {
		return nil, err
	}

	return lo.Map(subexpressions, func(subexpression dbSubexpression, _ int) model.Subexpression {
		return model.Subexpression(subexpression)
	}), nil
}

func (s *SubexpressionStorage) SubexpressionByDependableId(context context.Context, id int) (model.Subexpression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return model.Subexpression{}, err
	}
	defer conn.Close()

	var subexpression dbSubexpression

	arr := pq.Array([]int{id})

	if err := conn.GetContext(context,
		&subexpression,
		`SELECT id,
		 expression_id,
		 subexpression,
		 timeout FROM subexpressions WHERE depends_on = $1`, arr); err != nil {
		return model.Subexpression{}, err
	}

	return model.Subexpression(subexpression), nil
}

func (s *SubexpressionStorage) Delete(context context.Context, id int) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		context,
		`DELETE FROM subexpressions WHERE id = $1`,
		id,
	); err != nil {
		return err
	}

	return nil
}

func (s *SubexpressionStorage) CompleteSubexpression(context context.Context, id uuid.UUID) (float64, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var result float64

	if err := conn.GetContext(context, &result, `SELECT result FROM subexpressions WHERE expression_id = $1 ORDER BY id DESC LIMIT 1`, id); err != nil {
		return 0, err
	}

	return result, nil
}

type dbSubexpression struct {
	ID            int       `db:"id"`
	ExpressionId  uuid.UUID `db:"expression_id"`
	WorkerId      int       `db:"worker_id"`
	Subexpression string    `db:"subexpression"`
	IsTaken       bool      `db:"is_taken"`
	Timeout       float64   `db:"timeout"`
	DependsOn     []int     `db:"depends_on"`
	Result        float64   `db:"result"`
	Level         int       `db:"level"`
}
