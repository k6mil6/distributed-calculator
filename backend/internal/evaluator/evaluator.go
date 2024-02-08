package evaluator

func Evaluate(expression string, timeouts map[string]float64) (float64, error) {
	var timeout float64

	for i := range expression {
		for key, _ := range timeouts {
			if string(expression[i]) == key {
				if timeouts[key] > timeout {
					timeout = timeouts[key]
				}
			}
		}
	}

	return 0, nil
}
