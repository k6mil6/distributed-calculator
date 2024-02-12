package response

import (
	"github.com/google/uuid"
	"time"
)

type Response struct {
	Id         uuid.UUID     `json:"id"`
	Expression string        `json:"expression"`
	Timeout    time.Duration `json:"timeout"`
}
