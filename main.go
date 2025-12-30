package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var variables = make(map[string]int)

// Operator precedence
func precedence(op string) int {
	switch op {
	case "^":
		return 3
	case "*", "/":
		return 2
	case "+", "-":
		return 1
	}
	return 0
}

// Associativity: true if right-associative
func isRightAssociative(op string) bool {
	return op == "^"
}

// Check if valid identifier
func isValidIdentifier(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0
}

// Resolve value: number or variable
func resolveValue(token string) (int, error) {
	if isNumber(token) {
		return strconv.Atoi(token)
	}
	if isValidIdentifier(token) {
		val, ok := variables[token]
		if !ok {
			return 0, fmt.Errorf("Unknown variable")
		}
		return val, nil
	}
	return 0, fmt.Errorf("Invalid identifier")
}

// Check if number
func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// Convert infix to postfix using Shunting Yard
func infixToPostfix(expr string) ([]string, error) {
	// Normalize operators like +++ or ---
	expr = normalizeOperators(expr)

	tokens := tokenize(expr)
	output := []string{}
	stack := []string{}

	for _, token := range tokens {
		if isNumber(token) || isValidIdentifier(token) {
			output = append(output, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			foundLeft := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top == "(" {
					foundLeft = true
					break
				}
				output = append(output, top)
			}
			if !foundLeft {
				return nil, fmt.Errorf("Invalid expression")
			}
		} else if token == "+" || token == "-" || token == "*" || token == "/" || token == "^" {
			// Invalid sequences of * or /
			if strings.Contains(token, "**") || strings.Contains(token, "//") {
				return nil, fmt.Errorf("Invalid expression")
			}
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if precedence(top) > precedence(token) ||
					(precedence(top) == precedence(token) && !isRightAssociative(token)) {
					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, token)
		} else {
			return nil, fmt.Errorf("Invalid expression")
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top == "(" || top == ")" {
			return nil, fmt.Errorf("Invalid expression")
		}
		output = append(output, top)
	}
	return output, nil
}

// Tokenize expression (split into numbers, variables, operators, parentheses)
func tokenize(expr string) []string {
	// Add spaces around operators and parentheses
	replacer := strings.NewReplacer(
		"(", " ( ",
		")", " ) ",
		"+", " + ",
		"-", " - ",
		"*", " * ",
		"/", " / ",
		"^", " ^ ",
	)
	expr = replacer.Replace(expr)
	return strings.Fields(expr)
}

// Normalize sequences of + and -
func normalizeOperators(expr string) string {
	// Replace sequences of + with single +
	expr = strings.ReplaceAll(expr, "++", "+")
	// Replace sequences of -- with +
	for strings.Contains(expr, "--") {
		expr = strings.ReplaceAll(expr, "--", "+")
	}
	// Replace sequences of +- or -+ with -
	for strings.Contains(expr, "+-") {
		expr = strings.ReplaceAll(expr, "+-", "-")
	}
	for strings.Contains(expr, "-+") {
		expr = strings.ReplaceAll(expr, "-+", "-")
	}
	return expr
}

// Evaluate postfix expression
func evaluatePostfix(postfix []string) (int, error) {
	stack := []int{}
	for _, token := range postfix {
		if isNumber(token) || isValidIdentifier(token) {
			val, err := resolveValue(token)
			if err != nil {
				return 0, err
			}
			stack = append(stack, val)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("Invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var res int
			switch token {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("Division by zero")
				}
				res = a / b
			case "^":
				res = int(math.Pow(float64(a), float64(b)))
			default:
				return 0, fmt.Errorf("Invalid expression")
			}
			stack = append(stack, res)
		}
	}
	if len(stack) != 1 {
		return 0, fmt.Errorf("Invalid expression")
	}
	return stack[0], nil
}

// Handle assignment
func handleAssignment(line string) {
	parts := strings.Split(line, "=")
	if len(parts) != 2 {
		fmt.Println("Invalid assignment")
		return
	}
	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])

	if !isValidIdentifier(left) {
		fmt.Println("Invalid identifier")
		return
	}

	if isNumber(right) {
		val, _ := strconv.Atoi(right)
		variables[left] = val
		return
	}
	if isValidIdentifier(right) {
		val, ok := variables[right]
		if !ok {
			fmt.Println("Unknown variable")
			return
		}
		variables[left] = val
		return
	}
	fmt.Println("Invalid assignment")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())

		if line == "/exit" {
			fmt.Println("Bye!")
			break
		}
		if line == "/help" {
			fmt.Println("The program supports +, -, *, /, ^ and parentheses ().")
			fmt.Println("It also supports variables and unary minus.")
			continue
		}
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "/") {
			fmt.Println("Unknown command")
			continue
		}
		if strings.Contains(line, "=") {
			handleAssignment(line)
			continue
		}
		if isValidIdentifier(line) {
			val, ok := variables[line]
			if !ok {
				fmt.Println("Unknown variable")
			} else {
				fmt.Println(val)
			}
			continue
		}

		postfix, err := infixToPostfix(line)
		if err != nil {
			fmt.Println("Invalid expression")
			continue
		}
		result, err := evaluatePostfix(postfix)
		if err != nil {
			fmt.Println("Invalid expression")
			continue
		}
		fmt.Println(result)
	}
}
