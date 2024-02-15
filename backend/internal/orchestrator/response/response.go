package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK         = "OK"
	StatusInProgress = "In progress"
	StatusError      = "Error"
)

func InProgress() Response {
	return Response{
		Status: StatusInProgress,
	}
}

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
