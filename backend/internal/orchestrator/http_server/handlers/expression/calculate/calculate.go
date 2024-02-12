package calculate

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	resp "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/response"
	"github.com/k6mil6/distributed-calculator/backend/internal/storage"
	"github.com/k6mil6/distributed-calculator/backend/pkg/validation"
	"log/slog"
	"net/http"
)

type Request struct {
	Id         uuid.UUID      `json:"id"`
	Expression string         `json:"expression"`
	Timeouts   map[string]int `json:"timeouts"`
}

type Response struct {
	resp.Response
	Id uuid.UUID `json:"id"`
}

type ExpressionSaver interface {
	Save(context context.Context, expression model.Expression) error
}

func New(logger *slog.Logger, expressionSaver ExpressionSaver, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.expression.calculate.New"

		logger = logger.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("error decoding JSON request:", err)

			render.JSON(w, r, resp.Error("error decoding JSON request"))

			return
		}

		logger.Info("request body decoded", slog.Any("request", req))

		if !validation.IsMathExpressionValid(req.Expression) {
			logger.Error("invalid math expression")

			render.JSON(w, r, resp.Error("invalid math expression"))

			return
		}

		logger.Info("math expression is valid")

		err := expressionSaver.Save(context, model.Expression{
			ID:         req.Id,
			Expression: req.Expression,
			Timeouts:   req.Timeouts,
		})

		if err != nil {
			if errors.Is(err, storage.ErrExpressionInProgress) {
				responseOK(w, r, req.Id)

				return
			}
			logger.Error("error saving expression:", err)

			render.JSON(w, r, resp.Error("error saving expression"))

			return
		}
		logger.Info("expression saved successfully: ", req.Id)

		responseOK(w, r, req.Id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       id,
	})
}
