package fetcher

import (
	"fmt"
	shuntingYard "github.com/mgenware/go-shunting-yard"
	"testing"
)

func TestDivideIntoSubexpressions(t *testing.T) {
	infixTokens, err := shuntingYard.Scan("2^2")
	if err != nil {
		fmt.Println(err)
	}

	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		fmt.Println(err)
	}

	for _, token := range postfixTokens {
		fmt.Println(token.Type)
		fmt.Println(token.Value)
	}
}
