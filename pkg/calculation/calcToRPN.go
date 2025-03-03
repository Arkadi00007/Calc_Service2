package calculation

import (
	"fmt"
	"strings"
)

func PassSpace(exp string) string {
	return strings.ReplaceAll(exp, " ", "")
}

func IsDigit(s string) bool {
	return len(s) != 0 && s[0] >= '0' && s[0] <= '9'
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

func IsOperator(c uint8) bool {
	return c == '+' || c == '-' || c == '*' || c == '/'
}

func Calc(expression string) ([]string, error) {
	expression = PassSpace(expression)
	stack := make([]string, 0)
	que := make([]string, 0)

	if len(expression) == 0 {
		return nil, fmt.Errorf("empty expression")
	}

	if IsOperator(expression[len(expression)-1]) || IsOperator(expression[0]) || expression[0] == ')' {
		return nil, fmt.Errorf("invalid operator at position 0")
	}

	last := -1
	brackets := 0

	for i := 0; i < len(expression); i++ {
		v := expression[i]

		if !IsDigit(string(v)) && !IsOperator(v) && v != '(' && v != ')' {
			return nil, fmt.Errorf("invalid symbol '%c' at position %d", v, i)
		}

		if IsDigit(string(v)) {
			str := string(v)
			for i+1 < len(expression) && IsDigit(string(expression[i+1])) {
				i++
				str += string(expression[i])
			}
			que = append(que, str)
			last = i
		} else if IsOperator(v) {
			if last == -1 || IsOperator(expression[last]) {
				return nil, fmt.Errorf("invalid operator at position %d", i)
			}

			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(string(v)) {
				que = append(que, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, string(v))
			last = i
		} else if v == '(' {
			if last != -1 && IsDigit(string(expression[last])) {
				return nil, fmt.Errorf("missing operator before '(' at position %d", i)
			}
			stack = append(stack, string(v))
			brackets++
			last = i
		} else if v == ')' {
			if last == -1 || IsOperator(expression[last]) {
				return nil, fmt.Errorf("invalid ')' at position %d", i)
			}

			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top == "(" {
					break
				}
				que = append(que, top)
			}

			brackets--
			if brackets < 0 {
				return nil, fmt.Errorf("error: ')' without matching '(' at position %d", i)
			}
			last = i
		}
	}

	if brackets > 0 {
		return nil, fmt.Errorf("error: '(' without matching ')'")
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("error: '(' without matching ')'")
		}
		que = append(que, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return que, nil
}
