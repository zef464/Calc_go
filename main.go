package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, errors.New("empty expression")
	}

	rpn, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	return evaluateRPN(rpn)
}

func tokenize(expression string) []string {
	var tokens []string
	var current strings.Builder

	for _, char := range expression {
		if char == ' ' {
			continue
		}
		if isOperator(string(char)) || char == '(' || char == ')' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(char))
		} else if isIdentifier(char) || isDigit(char) || char == '.' {
			current.WriteRune(char)
		} else {
			return nil // Invalid character
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func isIdentifier(char rune) bool {
	return char >= 'a' && char <= 'z' // односимвольные идентификаторы
}

func isDigit(char rune) bool {
	return (char >= '0' && char <= '9') || char == '.'
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func shuntingYard(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	for _, token := range tokens {
		if isIdentifier(rune(token[0])) || isDigit(rune(token[0])) {
			output = append(output, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("mismatched parentheses")
			}
			stack = stack[:len(stack)-1] // Удаляем '('
		} else if isOperator(token) {
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(token) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else {
			return nil, errors.New("invalid token: " + token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if isIdentifier(rune(token[0])) {
			return 0, errors.New("identifiers are not supported in evaluation")
		} else if isDigit(rune(token[0])) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, errors.New("invalid RPN expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				stack = append(stack, a/b)
			default:
				return 0, errors.New("unknown operator: " + token)
			}
		} else {
			return 0, errors.New("invalid token in RPN: " + token)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid RPN expression")
	}

	return stack[0], nil
}

func main() {
	expression := "3 + 5 * (2 - 8)"
	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
