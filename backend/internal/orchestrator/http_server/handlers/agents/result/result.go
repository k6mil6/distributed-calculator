package result

import (
	"context"
	"github.com/go-chi/render"
	resp "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

type ExpressionResultSaver interface {
	SubexpressionIsDone(context context.Context, id int, result float64) error
}

type Request struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

type Response struct {
	resp.Response
	Id int `json:"id"`
}

func New(logger *slog.Logger, expressionResultSaver ExpressionResultSaver, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.agents.result.New"

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

		if err := expressionResultSaver.SubexpressionIsDone(context, req.Id, req.Result); err != nil {
			logger.Error("error saving expression result:", err)

			render.JSON(w, r, resp.Error("error saving expression result"))

			return
		}

		logger.Info("expression result saved", slog.Any("request", req))

		responseOK(w, r, req.Id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       id,
	})
}
