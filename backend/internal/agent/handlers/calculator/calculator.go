package calculator

type Request struct {
	Expression string
	RequestID  string
	Timeouts   map[string]float64
}

type Response struct {
	Result    float64
	Status    string
	RequestID string
}
