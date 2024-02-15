package subexpression_remaker

import (
	"fmt"
	"testing"
)

func TestSub(t *testing.T) {
	fmt.Println(Remake("{1} + {2}", 1, 2))
}
