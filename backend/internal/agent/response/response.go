package response

type Response struct {
	Id            int     `json:"id"`
	Subexpression string  `json:"subexpression"`
	Timeout       float64 `json:"timeout"`
	WorkerId      int     `json:"worker_id"`
}
