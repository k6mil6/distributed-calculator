package result

import (
	"context"
	"github.com/go-chi/render"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	resp "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/response"
	"github.com/k6mil6/distributed-calculator/backend/pkg/subexpression_remaker"
	"log/slog"
	"net/http"
)

type ExpressionResultSaver interface {
	SubexpressionIsDone(context context.Context, id int, result float64) error
	DoneSubexpressions(context context.Context) ([]model.Subexpression, error)
	SubexpressionByDependableId(context context.Context, id int) (model.Subexpression, error)
	Delete(context context.Context, id int) error
	Save(context context.Context, subExpression model.Subexpression) error
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

		doneSubexpressions, err := expressionResultSaver.DoneSubexpressions(context)
		if err != nil {
			logger.Error("error getting done subexpressions:", err)

			render.JSON(w, r, resp.Error("error getting done subexpressions"))

			return
		}

		for _, subexpression := range doneSubexpressions {

			// TODO: REMAKE THAT ONLY EXPRESSIONS THAT DEPEND ON ONE SUBEXPRESSION ARE BEING CALCULATED

			dependableSubexpression, err := expressionResultSaver.SubexpressionByDependableId(context, subexpression.ID)
			if err != nil {
				logger.Error("error getting subexpression:", err)

				render.JSON(w, r, resp.Error("error getting subexpression"))

				return
			}

			remadeSubexpression := subexpression_remaker.Remake(dependableSubexpression.Subexpression, subexpression.ID, subexpression.Result)
			dependableSubexpression.Subexpression = remadeSubexpression
			dependableSubexpression.DependsOn = nil

			if err := expressionResultSaver.Delete(context, dependableSubexpression.ID); err != nil {
				logger.Error("error deleting subexpression:", err)

				render.JSON(w, r, resp.Error("error deleting subexpression"))

				return
			}

			if err := expressionResultSaver.Save(context, dependableSubexpression); err != nil {
				logger.Error("error saving subexpression:", err)

				render.JSON(w, r, resp.Error("error saving subexpression"))

				return
			}
		}

		responseOK(w, r, req.Id)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       id,
	})
}
