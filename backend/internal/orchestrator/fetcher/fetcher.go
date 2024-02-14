package fetcher

import (
	"context"
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
		f.logger.Info("fetching expression", "expression_id", expression.ID)
	}
}

func divideIntoSubexpressions(expression model.Expression) ([]model.Subexpression, error) {
	infixTokens, err := shuntingYard.Scan(expression.Expression)
	if err != nil {
		panic(err)
	}

	fmt.Println("Infix Tokens:")
	fmt.Println(infixTokens)

	// Convert infix notation to postfix notation (RPN)
	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		panic(err)
	}

	fmt.Println("Postfix Tokens:")
	fmt.Println(postfixTokens)

	for _, token := range postfixTokens {
		fmt.Println(token)
	}

	return nil, nil
}
