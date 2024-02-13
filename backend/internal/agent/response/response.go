package response

type Response struct {
	Id            int    `json:"id"`
	Subexpression string `json:"subexpression"`
	Timeout       int64  `json:"timeout"`
}
