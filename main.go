package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var variables = make(map[string]int)

// normalizeOperator reduces sequences of + and - into a single + or -
func normalizeOperator(op string) string {
	countMinus := strings.Count(op, "-")
	if countMinus%2 == 0 {
		return "+"
	}
	return "-"
}

// isValidIdentifier checks if a variable name contains only Latin letters
func isValidIdentifier(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0
}

// resolveValue returns the integer value of a token (number or variable)
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

// isNumber checks if a string can be parsed as an integer
func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// evaluateExpression parses and evaluates an expression with + and - operators
func evaluateExpression(expr string) (int, error) {
	tokens := strings.Fields(expr)
	if len(tokens) == 0 {
		return 0, fmt.Errorf("empty expression")
	}

	result := 0
	sign := "+"
	for _, token := range tokens {
		// Operator sequence
		if strings.ContainsAny(token, "+-") && !isNumber(token) && !isValidIdentifier(token) {
			sign = normalizeOperator(token)
			continue
		}

		// Number or variable
		val, err := resolveValue(token)
		if err != nil {
			return 0, err
		}
		if sign == "+" {
			result += val
		} else {
			result -= val
		}
	}
	return result, nil
}

// handleAssignment processes variable assignment
func handleAssignment(line string) {
	parts := strings.Split(line, "=")
	if len(parts) != 2 {
		fmt.Println("Invalid assignment")
		return
	}
	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])

	// Validate left side
	if !isValidIdentifier(left) {
		fmt.Println("Invalid identifier")
		return
	}

	// Right side can be number or variable
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

		// Exit command
		if line == "/exit" {
			fmt.Println("Bye!")
			break
		}

		// Help command
		if line == "/help" {
			fmt.Println("The program calculates expressions with addition (+) and subtraction (-).")
			fmt.Println("It also supports storing values in variables.")
			continue
		}

		// Ignore empty lines
		if line == "" {
			continue
		}

		// Unknown command
		if strings.HasPrefix(line, "/") {
			fmt.Println("Unknown command")
			continue
		}

		// Assignment
		if strings.Contains(line, "=") {
			handleAssignment(line)
			continue
		}

		// Single variable
		if isValidIdentifier(line) {
			val, ok := variables[line]
			if !ok {
				fmt.Println("Unknown variable")
			} else {
				fmt.Println(val)
			}
			continue
		}

		// Expression
		if !isValidExpression(line) {
			fmt.Println("Invalid expression")
			continue
		}
		result, err := evaluateExpression(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(result)
	}
}

// isValidExpression checks if the expression contains only digits, spaces, +, -, or identifiers
func isValidExpression(expr string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9+\-\s]+$`)
	return re.MatchString(expr)
}
