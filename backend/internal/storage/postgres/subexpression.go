package postgres

import (
	"context"
	"database/sql"
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

	return lo.Map(subexpressions, func(subexpression dbSubexpression, _ int) model.Subexpression {
		return model.Subexpression(subexpression)
	}), nil
}

func (s *SubexpressionStorage) TakeSubexpression(ctx context.Context, id int) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var workerId int
	err = tx.QueryRowContext(ctx, "SELECT MAX(worker_id) FROM subexpressions").Scan(&workerId)
	if err != nil {
		return 0, err
	}
	workerId++

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE subexpressions SET is_taken = true, worker_id = $1 WHERE id = $2`,
		workerId,
		id,
	); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return workerId, nil
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

func (s *SubexpressionStorage) SubexpressionByDependableId(ctx context.Context, id int) ([]model.Subexpression, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `SELECT id, expression_id, subexpression, timeout, depends_on 
                                         FROM subexpressions WHERE $1 = ANY(depends_on)`, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var subexpressions []model.Subexpression
	for rows.Next() {
		var subexpr model.Subexpression
		var dependsOn pq.Int64Array // Use pq.Int64Array for PostgreSQL integer array
		if err := rows.Scan(&subexpr.ID, &subexpr.ExpressionId, &subexpr.Subexpression, &subexpr.Timeout, &dependsOn); err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Convert pq.Int64Array to []int
		subexpr.DependsOn = make([]int, len(dependsOn))
		for i, v := range dependsOn {
			subexpr.DependsOn[i] = int(v)
		}

		// Append the constructed model.Subexpression to the slice
		subexpressions = append(subexpressions, subexpr)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return subexpressions, nil
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

func (s *SubexpressionStorage) LastWorkerId(context context.Context) (int, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var workerId sql.NullInt64

	if err := conn.GetContext(context, &workerId, `SELECT worker_id FROM subexpressions ORDER BY worker_id LIMIT 1`); err != nil {
		fmt.Println(err)
		return 0, err
	}

	if !workerId.Valid {
		fmt.Println(workerId.Valid)
		return 0, nil
	}

	fmt.Println(workerId.Int64)

	return int(workerId.Int64), nil
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
