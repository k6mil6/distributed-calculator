package models

import "time"

type Expression struct {
	ID         int64
	Expression string
	CreatedAt  time.Time
	IsTaken    bool
	IsDone     bool
}

type Subexpression struct {
	ID            int64
	ExpressionId  int64
	WorkerId      int64
	Subexpression string
	IsTaken       bool
	IsDone        bool
	Timeout       time.Duration
	Result        float64
}
