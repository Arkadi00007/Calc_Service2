package calculation

import (
	"reflect"
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		expr     string
		expected []string
		hasError bool
	}{

		{"1+2", []string{"1", "2", "+"}, false},
		{"3-4", []string{"3", "4", "-"}, false},
		{"5*6", []string{"5", "6", "*"}, false},
		{"8/2", []string{"8", "2", "/"}, false},
		{"(1+2)", []string{"1", "2", "+"}, false},
		{"(3+4)*5", []string{"3", "4", "+", "5", "*"}, false},
		{"6/(2+1)", []string{"6", "2", "1", "+", "/"}, false},
		{"(8-3)*4", []string{"8", "3", "-", "4", "*"}, false},
		{"2+(3*4)", []string{"2", "3", "4", "*", "+"}, false},
		{"(1+2)*(3+4)", []string{"1", "2", "+", "3", "4", "+", "*"}, false},

		{"(1+2", nil, true},   // Нет закрывающей скобки
		{"1+2)", nil, true},   // Нет открывающей скобки
		{"((3+4)", nil, true}, // Лишняя открывающая скобка
		{"(5-2))", nil, true}, // Лишняя закрывающая скобка
		{"(()", nil, true},    // Пустая скобка
		{")1+2(", nil, true},  // Закрывающая перед открывающей

		{"+1+2", nil, true}, // Оператор в начале
		{"1+2+", nil, true}, // Оператор в конце
		{"1++2", nil, true}, // Два подряд идущих оператора
		{"1+*2", nil, true}, // Несовместимые операторы

		{"1*/2", nil, true}, // Неверный оператор между числами

		{"a+b", nil, true},  // Буквы
		{"1+2@", nil, true}, // Спецсимволы

		{"1+2#", nil, true}, // Неизвестные символы

		{"", nil, true}, // Пустой ввод
	}

	for _, test := range tests {
		result, err := Calc(test.expr)

		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", test.expr)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.expr, err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("For input '%s': expected %v, got %v", test.expr, test.expected, result)
			}
		}
	}
}
