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
	UnkownOperator        error = errors.New("unkown operator") //TODO: Убрать неиспользуемую ошибку
	ExpressionEmpty       error = errors.New("expression empty")
	OperationWithoutValue error = errors.New("operation dont have a value")
	BracketsNotFound      error = errors.New("not found opened or closed bracket")
	ParseError            error = errors.New("parse error")
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
		op_bracket, cl_bracket := strings.IndexRune(ex, '('), strings.IndexRune(ex, ')')

		//TODO: Исправить не нахождение последней скобки
		//* Пример ошибки: GetExample("(10+30+(20+1*10))+1") got "(10+30+(20+1*10)", but expected "(10+30+(20+1*10))"

		if (op_bracket == -1 && cl_bracket != -1) || (op_bracket != -1 && cl_bracket == -1) {
			return Example{}, BracketsNotFound
		}

		return Example{String: ex[op_bracket : cl_bracket+1]}, nil
	}

	operator_idx = strings.IndexAny(ex, "*/")
	if operator_idx == -1 {
		operator_idx = strings.IndexAny(ex, "+-")
		if operator_idx == -1 {
			value, err := strconv.ParseFloat(ex, 64)
			if err != nil {
				return Example{}, ParseError
			}
			return Example{Second_value: value, Operation: Equals}, nil
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

// TODO: Удалить нахер эту поеботу
func GetExample(example string) (string, int, Example, error) {
	var ex Example
	var local_ex string = example

	var begin, end int = -1, -1
	if strings.ContainsRune(local_ex, '(') {
		for i, rn := range local_ex {
			if rn == '(' {
				begin = i
				continue
			} else if rn == ')' {
				end = i
				break
			}
		}

		if (begin == -1 && end != -1) || (begin != -1 && end == -1) {
			return "", 0, Example{}, BracketsNotFound
		}

		local_ex = local_ex[begin : end+1]
	}

	var actionIdx int
	if op := "*/"; strings.ContainsAny(local_ex, op) {
		actionIdx = strings.IndexAny(local_ex, op)
	} else if op := "+-"; strings.ContainsAny(local_ex, op) {
		actionIdx = strings.IndexAny(local_ex, op)
	} else if strings.ContainsAny(local_ex, "()") {
		value, err := strconv.ParseFloat(local_ex[1:len(local_ex)-1], 64)
		if err != nil {
			return "", 0, Example{}, ParseError
		}
		return local_ex[:], strings.IndexRune(example, rune(local_ex[0])), Example{First_value: value, Second_value: 52, Operation: Equals}, nil //52 - по рофлу, чтобы при калькулировании не возникала ошибка. Крч костыль
	} else {
		value, err := strconv.ParseFloat(local_ex, 64)
		if err != nil {
			return "", 0, Example{}, ParseError
		}
		return "end", 0, Example{First_value: value, Second_value: 52, Operation: Equals}, nil
	}

	if actionIdx == 0 || actionIdx == len(local_ex)-1 {
		return "", 0, Example{}, OperationWithoutValue
	}

	ex.Operation = Operator(local_ex[actionIdx])

	//Нахождение концов двух чисел
	var exampleLen = len(local_ex)
	if actionIdx == 0 || actionIdx == exampleLen-1 {

		//TODO: Изменить вывод ошибки на вывод с использованием переменной

		return "", 0, Example{}, errors.New("action in first or lst place")
	}

	var err error
	for i := actionIdx - 1; i >= 0; i-- {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			ex.First_value, err = strconv.ParseFloat(local_ex[i+1:actionIdx], 64)
			if err != nil {
				return "", 0, Example{}, ParseError
			}
			begin = i + 1
			break
		} else if i == 0 {
			ex.First_value, err = strconv.ParseFloat(local_ex[i:actionIdx], 64)
			if err != nil {
				return "", 0, Example{}, ParseError
			}
			begin = i
			break
		}
	}

	for i := actionIdx + 1; i < exampleLen; i++ {
		if strings.ContainsRune("+-/*()", rune(local_ex[i])) {
			ex.Second_value, err = strconv.ParseFloat(local_ex[actionIdx+1:i], 64)
			if err != nil {
				return "", 0, Example{}, ParseError
			}
			end = i
			break
		} else if i+1 == exampleLen {
			ex.Second_value, err = strconv.ParseFloat(local_ex[actionIdx+1:i+1], 64)
			if err != nil {
				return "", 0, Example{}, ParseError
			}
			end = exampleLen
			break
		}
	}

	return local_ex[begin:end], strings.IndexRune(example, rune(local_ex[0])), ex, nil
}

func SolveExample(ex Example) (float64, error) {
	if ex.Second_value == 0 {
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
	return 0, UnkownOperator //TODO: Убрать неиспользуемую ошибку
}

func EraseExample(example, erase_ex string, pri_idx int, answ float64) string {
	return example[:pri_idx] + strings.Replace(example[pri_idx:], erase_ex, fmt.Sprintf("%f", answ), 1)
}

func Calc(expression string) (result float64, err error) {
	if expression == "" {
		return 0, ExpressionEmpty
	}

	for {
		ex_str, pri_idx, example, err := GetExample(expression)
		if err != nil {
			return 0, err
		}

		result, _ = SolveExample(example)

		if ex_str == "end" {
			break
		}

		expression = EraseExample(expression, ex_str, pri_idx, result)
	}
	return
}

func main() {}
