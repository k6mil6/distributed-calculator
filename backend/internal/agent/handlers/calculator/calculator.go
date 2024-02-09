package calculator

import (
	"errors"
	"github.com/go-chi/render"
	resp "github.com/k6mil6/distributed-calculator/backend/lib/api/response"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	ID         string  `json:"id,required"`
	Expression string  `json:"expression,required"`
	Timeout    float64 `json:"timeouts,required"`
}

type Response struct {
	resp.Response
	Result    float64 `json:"result,required"`
	RequestID string  `json:"request_id,required"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("request body is empty"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}
		if err != nil {
			log.Error("request body is invalid")

			render.JSON(w, r, resp.Error("request body is invalid"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		log.Info("request body is valid")

	}
}
