package set_timeouts

import (
	"context"
	"github.com/go-chi/render"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	resp "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/response"
	"github.com/k6mil6/distributed-calculator/backend/internal/timeout"
	"log/slog"
	"net/http"
)

type Request struct {
	Timeouts timeout.Timeout `json:"timeouts"`
}

type TimeoutsSetter interface {
	Save(context context.Context, timeouts model.Timeouts) error
}

func New(logger *slog.Logger, timeoutsSetter TimeoutsSetter, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.timeouts.set_timeouts.New"

		logger = logger.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("error decoding JSON request:", err)

			render.JSON(w, r, resp.Error("error decoding JSON request"))

			return
		}

		if err := timeoutsSetter.Save(context, model.Timeouts{
			Timeouts: req.Timeouts,
		}); err != nil {
			logger.Error("error setting timeouts:", err)

			render.JSON(w, r, resp.Error("error setting timeouts"))

			return

		}

		render.JSON(w, r, resp.OK())
	}
}
