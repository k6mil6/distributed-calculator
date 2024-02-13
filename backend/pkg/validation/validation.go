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
	if !divisionByZero(expression) {
		return false
	}
	return true
}

func divisionByZero(expression string) bool {
	expression = strings.ReplaceAll(expression, " ", "")

	for i := range expression {
		if expression[i] == '/' {
			if i+1 <= len(expression)-1 {
				return false
			}
		}
	}

	return true
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

	for _, char := range expression {
		if char == ' ' {
			continue
		}
		if char == '=' {
			if lastWasOperator {
				return false
			}
			lastWasOperator = true
		} else if isOperator(char) {
			if lastWasOperator && lastWasMinus {
				return false
			} else if lastWasOperator && char == '-' {
				lastWasMinus = true
			} else if lastWasOperator {
				return false
			} else {
				lastWasMinus = false
				lastWasOperator = true
			}
		} else if char != '(' && char != ')' {
			lastWasMinus = false
			if !lastWasOperator {
				return false
			} else {
				lastWasOperator = false
				if unicode.IsLetter(char) {
					for i, c := range expression {
						if unicode.IsLetter(c) {
							continue
						}
						if i+1 == len(expression) || !unicode.IsLetter(rune(expression[i+1])) {
							break
						}
					}
				} else if unicode.IsNumber(char) {
					comma := false
					for i, c := range expression {
						if unicode.IsNumber(c) {
							continue
						}
						if c == ',' {
							if comma {
								return false
							}
							comma = true
						}
						if i+1 == len(expression) || !unicode.IsNumber(rune(expression[i+1])) {
							break
						}
					}
				}
			}
		}
	}
	return !lastWasOperator
}

func isOperator(character rune) bool {
	switch character {
	case '+', '-', '*', '^', '/':
		return true
	}
	return false
}
