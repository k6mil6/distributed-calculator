package model

import (
	"github.com/google/uuid"
	"time"
)

type Expression struct {
	ID         uuid.UUID
	Expression string
	CreatedAt  time.Time
	IsTaken    bool
	IsDone     bool
}

type Subexpression struct {
	ID            int
	ExpressionId  uuid.UUID
	WorkerId      int
	Subexpression string
	IsTaken       bool
	IsDone        bool
	Timeout       time.Duration
	Result        float64
}
