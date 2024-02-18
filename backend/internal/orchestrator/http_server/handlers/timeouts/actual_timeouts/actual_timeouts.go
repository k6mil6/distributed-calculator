package actual_timeouts

import (
	"context"
	"github.com/go-chi/render"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	resp "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

type TimeoutsGetter interface {
	GetActualTimeouts(context context.Context) (model.Timeouts, error)
}

type Response struct {
	Timeouts model.Timeouts `json:"timeouts"`
}

func New(logger *slog.Logger, timeoutsGetter TimeoutsGetter, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.timeouts.actual_timeouts.New"

		logger = logger.With(
			slog.String("op", op),
		)

		timeouts, err := timeoutsGetter.GetActualTimeouts(context)
		if err != nil {
			logger.Error("error getting actual timeouts:", err)

			render.JSON(w, r, resp.Error("error getting actual timeouts"))

			return
		}

		render.JSON(w, r, Response{
			Timeouts: timeouts,
		})
	}
}
