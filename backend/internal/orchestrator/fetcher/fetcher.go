package fetcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	shuntingYard "github.com/mgenware/go-shunting-yard"
	"log/slog"
	"time"
)

type ExpressionGetter interface {
	NonTakenExpressions(context context.Context) ([]model.Expression, error)
	TakeExpression(context context.Context, id uuid.UUID) error
}

type SubexpressionSaver interface {
	Save(context context.Context, subExpression model.Subexpression) error
	LastSubexpression(context context.Context) (int, error)
}

type Fetcher struct {
	expressionGetter   ExpressionGetter
	subexpressionSaver SubexpressionSaver

	fetchInterval time.Duration
	logger        *slog.Logger
}

func New(expressionGetter ExpressionGetter, subexpressionSaver SubexpressionSaver, fetchInterval time.Duration, logger *slog.Logger) *Fetcher {
	return &Fetcher{
		expressionGetter:   expressionGetter,
		subexpressionSaver: subexpressionSaver,
		fetchInterval:      fetchInterval,
		logger:             logger,
	}
}

func (f *Fetcher) Start(context context.Context) {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	f.logger.Info("fetcher started")

	for {
		select {
		case <-ticker.C:
			f.Fetch(context)
		case <-context.Done():
			return
		}
	}
}

func (f *Fetcher) Fetch(context context.Context) {
	expressions, err := f.expressionGetter.NonTakenExpressions(context)
	if err != nil {
		f.logger.Error("failed to fetch expressions", err)
		return
	}

	for _, expression := range expressions {
		if err := f.expressionGetter.TakeExpression(context, expression.ID); err != nil {
			f.logger.Error("failed to take expression", err)
			return
		}
		f.logger.Info("fetching expression", "expression_id", expression.ID)

		subexpressions, err := f.divideIntoSubexpressions(context, expression)
		if err != nil {
			f.logger.Error("failed to divide expression into subexpressions", err)
			return
		}
		for _, subexpression := range subexpressions {
			f.logger.Info("saving subexpression", "subexpression_id", subexpression.ID, "expression", subexpression.Subexpression)
			err := f.subexpressionSaver.Save(context, subexpression)
			if err != nil {
				f.logger.Error("failed to save subexpression", err)
				return
			}
		}

	}
}

func (f *Fetcher) divideIntoSubexpressions(context context.Context, expression model.Expression) ([]model.Subexpression, error) {
	infixTokens, err := shuntingYard.Scan(expression.Expression)
	if err != nil {
		return nil, err
	}

	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		return nil, err
	}

	var subexpressions []model.Subexpression

	var stack []struct {
		expr string
		id   int
	}
	subExprID, err := f.subexpressionSaver.LastSubexpression(context)

	for _, token := range postfixTokens {
		if token.Type == 1 {
			numberStr := fmt.Sprintf("%d", token.Value)
			stack = append(stack, struct {
				expr string
				id   int
			}{expr: numberStr, id: -1})
		} else {
			if len(stack) < 2 {
				return nil, errors.New("not enough operands for the operator")
			}

			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			subExprID++
			subExpression := fmt.Sprintf("%s %s %s", left.expr, token.Value, right.expr)
			placeholder := fmt.Sprintf("{%d}", subExprID)

			subExpr := model.Subexpression{
				ExpressionId:  expression.ID,
				ID:            subExprID,
				Subexpression: subExpression,
				Timeout:       expression.Timeouts[token.Value.(string)].(float64),
				DependsOn:     []int{},
			}

			if left.id != -1 {
				subExpr.DependsOn = append(subExpr.DependsOn, left.id)
			}
			if right.id != -1 {
				subExpr.DependsOn = append(subExpr.DependsOn, right.id)
			}

			stack = append(stack, struct {
				expr string
				id   int
			}{expr: placeholder, id: subExprID})

			if len(subExpr.DependsOn) == 0 {
				subExpr.DependsOn = nil
			}

			subexpressions = append(subexpressions, subExpr)
		}
	}

	return subexpressions, nil
}
