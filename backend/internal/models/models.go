package models

import "time"

type Expression struct {
	ID         int64
	Expression string
	CreatedAt  time.Time
	IsTaken    bool
}
