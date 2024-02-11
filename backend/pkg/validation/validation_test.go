package validation

import "testing"

func TestIsMathExpressionValid(t *testing.T) {
	if !IsMathExpressionValid("(3+5)*2") {
		t.Error("Expected true, got false")
	}

	if IsMathExpressionValid("(3+5)*2)") {
		t.Error("Expected false, got true")
	}

	if IsMathExpressionValid("(+3+5)*2") {
		t.Error("Expected false, got true")
	}

	if !IsMathExpressionValid("((3+5)*2)") {
		t.Error("Expected true, got false")
	}

	if IsMathExpressionValid(")(3+5)*2)") {
		t.Error("Expected false, got true")
	}

	if IsMathExpressionValid(")(3++5)*2)") {
		t.Error("Expected false, got true")
	}
}
