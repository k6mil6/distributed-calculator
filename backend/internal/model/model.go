package model

import (
	"github.com/google/uuid"
	"time"
)

type Expression struct {
	ID         uuid.UUID
	Expression string
	CreatedAt  time.Time
	Timeouts   map[string]int
	IsTaken    bool
	IsDone     bool
	Result     float64
}

type Subexpression struct {
	ID            int
	ExpressionId  uuid.UUID
	WorkerId      int
	Subexpression string
	IsTaken       bool
	Timeout       int64
	DependsOn     []int
	Result        float64
}
