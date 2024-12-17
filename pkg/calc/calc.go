package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//TODO: Изменить вывод ошибок strconv.Parse... в вывод ошибок через переменную

type Operator rune

const (
	Plus     Operator = '+'
	Minus    Operator = '-'
	Multiply Operator = '*'
	Division Operator = '/'
	Equals   Operator = '='
)

var (
	DivideByZero          error = errors.New("divide by zero")
	ExpressionEmpty       error = errors.New("expression empty")
	OperationWithoutValue error = errors.New("operation dont have a value")
	BracketsNotFound      error = errors.New("not found opened or closed bracket")
	ParseError            error = errors.New("parse error")
	UnexpextedError       error = errors.New("unexpected error")
)

type Example struct {
	First_value  float64
	Second_value float64
	Operation    Operator
	String       string
}

/*
Return indexs first and second number value

If index first and second value = operator index - return error (OperationWithoutValue)
*/
func getValuesIdx(example_str string, operator_idx int) (first_value_idx int, second_value_idx int, err error) {
	first_value_idx, second_value_idx = operator_idx-1, operator_idx+1

	if first_value_idx < 0 || second_value_idx > len(example_str) {
		fmt.Println(example_str)
		return 0, 0, OperationWithoutValue
	}

	// Get first value index
	for ; first_value_idx != 0; first_value_idx-- {
		ex_rune := example_str[first_value_idx]
		if ex_rune == '+' || ex_rune == '-' || ex_rune == '*' || ex_rune == '/' {
			first_value_idx++
			break
		}
	}

	if operator_idx == first_value_idx {
		return 0, 0, OperationWithoutValue
	}

	// Get second value index
	for ; second_value_idx < len(example_str)-1; second_value_idx++ {
		ex_rune := example_str[second_value_idx]
		if ex_rune == '+' || ex_rune == '-' || ex_rune == '*' || ex_rune == '/' {
			second_value_idx--
			break
		}
	}

	if operator_idx == second_value_idx {
		return 0, 0, OperationWithoutValue
	}
	return
}

func GetExampleNew(ex string) (Example, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("GetExampleNew():", err)
		}
	}()

	var operator_idx int

	if strings.ContainsRune(ex, '(') || strings.ContainsRune(ex, ')') {
		opened_bracket := strings.IndexRune(ex, '(')
		closed_bracket := strings.IndexRune(ex, ')')

		if (opened_bracket == -1 && closed_bracket != -1) || (opened_bracket != -1 && closed_bracket == -1) {
			return Example{}, BracketsNotFound
		}

		var count int
		for i := opened_bracket + 1; i < len(ex); i++ {
			if rune(ex[i]) == ')' && count == 0 {
				closed_bracket = i
				break
			}

			symbol := rune(ex[i])
			if symbol == '(' {
				count++
			} else if symbol == ')' {
				count--
			}
		}

		if closed_bracket == opened_bracket {
			return Example{}, BracketsNotFound
		}

		return Example{String: ex[opened_bracket : closed_bracket+1]}, nil
	}

	operator_idx = strings.IndexAny(ex, "*/")
	if operator_idx == -1 {
		operator_idx = strings.IndexAny(ex, "+-")
		if operator_idx == -1 {
			value, err := strconv.ParseFloat(ex, 64)
			if err != nil {
				return Example{}, ParseError
			}
			fmt.Println("end", value)
			return Example{First_value: value, Operation: Equals}, nil
		}
	}

	first_value_idx, second_value_idx, err := getValuesIdx(ex, operator_idx)
	if err != nil {
		return Example{}, err
	}

	first_value, err := strconv.ParseFloat(ex[first_value_idx:operator_idx], 64)
	if err != nil {
		return Example{}, ParseError
	}

	second_value, err := strconv.ParseFloat(ex[operator_idx+1:second_value_idx+1], 64)
	if err != nil {
		return Example{}, ParseError
	}

	new_ex := Example{
		First_value:  first_value,
		Second_value: second_value,
		Operation:    Operator(ex[operator_idx]),
		String:       ex[first_value_idx : second_value_idx+1],
	}
	return new_ex, nil
}

func SolveExample(ex Example) (float64, error) {
	if ex.Second_value == 0 && ex.Operation == Division {
		return 0, DivideByZero
	}

	switch ex.Operation {
	case Plus:
		return ex.First_value + ex.Second_value, nil
	case Minus:
		return ex.First_value - ex.Second_value, nil
	case Multiply:
		return ex.First_value * ex.Second_value, nil
	case Division:
		return ex.First_value / ex.Second_value, nil
	case Equals:
		return ex.First_value, nil
	}
	return 0, UnexpextedError
}

func EraseExample(example, erase_ex string, pri_idx int, answ float64) string {
	return example[:pri_idx] + strings.Replace(example[pri_idx:], erase_ex, fmt.Sprintf("%f", answ), 1)
}

func Calc(expression string) (float64, error) {
	if expression == "" {
		return 0, ExpressionEmpty
	}

	var ex_with_brackets string
	ex_to_solve := expression
	local_ex := ex_to_solve

	for {
		fmt.Println("[LOG] ex_with_brackets:", ex_with_brackets, "local_ex in brackets:", local_ex)
		ex, err := GetExampleNew(local_ex)
		if err != nil {
			return 0, err
		}

		local_ex = ex.String

		if strings.ContainsAny(local_ex, "()") {
			ex_with_brackets = local_ex
			ex_to_solve = ex_with_brackets[1 : len(ex_with_brackets)-1]
			local_ex = ex_to_solve
			fmt.Println("[LOG] ex_with_brackets:", ex_with_brackets, "ex_to_solve in brackets:", ex_to_solve)
			continue
		}

		answ, err := SolveExample(ex)
		if err != nil {
			return 0, err
		}

		if ex.Operation == Equals && ex_with_brackets == "" {
			fmt.Println("ret", ex.Operation)
			return answ, nil
		}

		var idx_str int
		if ex_with_brackets != "" {
			idx_str = strings.Index(expression, ex_with_brackets)
		} else {
			idx_str = strings.Index(expression, local_ex)
		}

		if idx_str == -1 {
			return 0, UnexpextedError
		}

		expression = EraseExample(expression, local_ex, idx_str, answ)
		local_ex = expression
	}
}
