package response

import "time"

type Response struct {
	Id         int64         `json:"id"`
	Expression string        `json:"expression"`
	Timeout    time.Duration `json:"timeout"`
}
