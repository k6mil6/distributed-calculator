package free_expressions

import (
	"context"
	"github.com/go-chi/render"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	resp "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/response"
	"log/slog"
	"net/http"
	"time"
)

type SubexpressionGetter interface {
	NonTakenSubexpressions(context context.Context) ([]model.Subexpression, error)
	TakeSubexpression(context context.Context, id int) error
}

type Response struct {
	resp.Response
	Id            int           `json:"id"`
	Subexpression string        `json:"Subexpression"`
	Timeout       time.Duration `json:"timeout"`
}

func New(logger *slog.Logger, getter SubexpressionGetter, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.agents.free_expression.New"

		logger = logger.With(
			slog.String("op", op),
		)

		latestSubexpressions, err := getter.NonTakenSubexpressions(context)
		if err != nil {
			logger.Error("no available expressions", err)

			render.JSON(w, r, resp.Error("no available expressions"))

			return
		}

		subexpression := latestSubexpressions[0]

		if err := getter.TakeSubexpression(context, subexpression.ID); err != nil {
			logger.Error("error taking subexpression", err)

			render.JSON(w, r, resp.Error("error taking subexpression"))

			return
		}

		logger.Info("subexpression taken", slog.Int("id", subexpression.ID))

		responseOK(w, r, subexpression)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, subexpression model.Subexpression) {
	render.JSON(w, r, Response{
		Response:      resp.OK(),
		Id:            subexpression.ID,
		Subexpression: subexpression.Subexpression,
		Timeout:       subexpression.Timeout,
	})
}
