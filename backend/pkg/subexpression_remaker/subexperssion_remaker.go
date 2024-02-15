package subexpression_remaker

import (
	"fmt"
	"strings"
)

func Remake(subexpression string, id int, result float64) string {
	return strings.ReplaceAll(subexpression, fmt.Sprintf("{%v}", id), fmt.Sprint(result))
}
