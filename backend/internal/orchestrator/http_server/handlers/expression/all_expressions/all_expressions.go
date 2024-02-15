package all_expressions

import (
	"context"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	resp "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

type ExpressionsSelector interface {
	AllExpressions(context context.Context) ([]model.Expression, error)
	UpdateResult(context context.Context, id uuid.UUID, result float64) error
}

type SubexpressionCompleter interface {
	CompleteSubexpression(context context.Context, id uuid.UUID) (float64, error)
}

func New(logger *slog.Logger, expressionsSelector ExpressionsSelector, subexpressionCompleter SubexpressionCompleter, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.expression.all_expressions.New"

		logger = logger.With(
			slog.String("op", op),
		)

		expressions, err := expressionsSelector.AllExpressions(context)
		if err != nil {
			logger.Error("error getting all expressions:", err)

			render.JSON(w, r, resp.Error("error getting all expressions"))

			return
		}

		for _, expression := range expressions {
			if expression.Result != 0 {
				continue
			}

			result, err := subexpressionCompleter.CompleteSubexpression(context, expression.ID)
			if err != nil {
				logger.Error("expression is not done", err)

				continue
			}

			expression.Result = result

			if err := expressionsSelector.UpdateResult(context, expression.ID, result); err != nil {
				logger.Error("error updating expression result:", err)

				render.JSON(w, r, resp.Error("error updating expression result"))

				return
			}

			logger.Info("expression result updated", slog.Any("expression", expression))
		}

		logger.Info("expressions retrieved", slog.Any("expressions", expressions))

		render.JSON(w, r, expressions)
	}
}
