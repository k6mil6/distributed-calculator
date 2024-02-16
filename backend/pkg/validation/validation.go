package validation

import (
	"strings"
	"unicode"
)

func IsMathExpressionValid(expression string) bool {
	if !checkParentheses(expression) {
		return false
	}
	if !checkOperationOrder(expression) {
		return false
	}
	if divisionByZero(expression) {
		return false
	}
	return true
}

func divisionByZero(expression string) bool {
	expression = strings.ReplaceAll(expression, " ", "")

	for i := 0; i < len(expression)-1; i++ {
		if expression[i] == '/' && expression[i+1] == '0' {
			return true
		}
	}

	return false
}
func checkParentheses(expression string) bool {
	parenthesesOpen := 0
	for _, char := range expression {
		if char == ')' {
			parenthesesOpen--
		} else if char == '(' {
			parenthesesOpen++
		} else if char == '=' {
			if parenthesesOpen != 0 {
				return false
			}
		}
		if parenthesesOpen < 0 {
			return false
		}
	}
	return parenthesesOpen == 0
}

func checkOperationOrder(expression string) bool {
	lastWasOperator := true
	lastWasMinus := false
	lastWasNumber := false

	for i, char := range expression {
		if char == ' ' {
			continue
		}
		if isOperator(char) {
			if lastWasOperator && !lastWasMinus {
				return false
			}
			lastWasOperator = true
			lastWasMinus = char == '-'
			lastWasNumber = false
		} else if unicode.IsNumber(char) || (char == '.' && lastWasNumber) {
			if i > 0 && (unicode.IsNumber(rune(expression[i-1])) || expression[i-1] == '.') {
				continue
			}
			lastWasOperator = false
			lastWasMinus = false
			lastWasNumber = true
		} else if char != '(' && char != ')' {
			return false
		}
	}
	return !lastWasOperator
}

func isOperator(c rune) bool {
	switch c {
	case '+', '-', '*', '/':
		return true
	}
	return false
}
