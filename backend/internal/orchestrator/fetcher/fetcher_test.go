package fetcher

import (
	"github.com/k6mil6/distributed-calculator/backend/internal/model"
	"testing"
)

func TestDivideIntoSubexpressions(t *testing.T) {
	// Testing for valid expression
	validExpression := model.Expression{Expression: "3 + 4 * 2 / ( 1 - 5 ) ^ 2"}

	_, err := divideIntoSubexpressions(validExpression)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

}
