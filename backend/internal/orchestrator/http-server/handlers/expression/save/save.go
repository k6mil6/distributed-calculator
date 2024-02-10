package save

import (
	"context"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	resp "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/response"
	"log"
	"net/http"
)

type Request struct {
	Id         uuid.UUID `json:"id"`
	Expression string    `json:"expression"`
}

type Response struct {
	resp.Response
	Id uuid.UUID `json:"id"`
}

type ExpressionSaver interface {
	Save(context context.Context, expression model.Expression) error
}

func New(logger *log.Logger, expressionSaver ExpressionSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
