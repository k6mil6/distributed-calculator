package model

import (
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/backend/internal/timeout"
	"time"
)

type Expression struct {
	ID         uuid.UUID
	Expression string
	CreatedAt  time.Time
	Timeouts   timeout.Timeout
	IsTaken    bool
	Result     float64
}

type Subexpression struct {
	ID            int
	ExpressionId  uuid.UUID
	WorkerId      int
	Subexpression string
	IsTaken       bool
	Timeout       float64
	DependsOn     []int
	Result        float64
	Level         int
}
