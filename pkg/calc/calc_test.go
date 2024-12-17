package calc_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

func TestSolveExample(t *testing.T) {
	cases := []struct {
		name         string
		example      calc.Example
		expected     float64
		expected_err error
	}{
		{
			name:         "1 + 1",
			example:      calc.Example{First_value: 1, Second_value: 1, Operation: calc.Plus},
			expected:     2,
			expected_err: nil,
		},
		{
			name:         "1 * 1",
			example:      calc.Example{First_value: 1, Second_value: 1, Operation: calc.Multiply},
			expected:     1,
			expected_err: nil,
		},
		{
			name:         "divide by zero",
			example:      calc.Example{First_value: 0, Second_value: 0, Operation: calc.Division},
			expected:     0,
			expected_err: calc.DivideByZero,
		},
		{
			name:         "123 + 10",
			example:      calc.Example{First_value: 123, Second_value: 10, Operation: calc.Plus},
			expected:     133,
			expected_err: nil,
		},
		{
			name:         "equal(1)",
			example:      calc.Example{First_value: 1, Second_value: 52, Operation: calc.Equals},
			expected:     1,
			expected_err: nil,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := calc.SolveExample(test.example)
			if got != test.expected || !errors.Is(err, test.expected_err) {
				t.Errorf("SolveExample(%#v) = (%f, %q), but expected: (%f, %q)", test.example, got, err, test.expected, test.expected_err)
			}
		})
	}
}

func TestGetExample(t *testing.T) {
	cases := []struct {
		name         string
		example      string
		expected_str string
		expected_err error
	}{
		{
			name:         "1+1",
			example:      "1+1",
			expected_str: "1+1",
		},
		{
			name:         "1+145",
			example:      "1+145",
			expected_str: "1+145",
		},
		{
			name:         "10+10+10+10",
			example:      "10+10",
			expected_str: "10+10",
		},
		{
			name:         "1+145*10",
			example:      "1+145*10",
			expected_str: "145*10",
		},
		{
			name:         "(1)",
			example:      "(1)",
			expected_str: "(1)",
		},
		{
			name:         "(112+1)+145*10",
			example:      "(112+1)+145*10",
			expected_str: "(112+1)",
		},
		{
			name:         "(10+30+(20+1*10))+1",
			example:      "(10+30+(20+1*10))+1",
			expected_str: "(10+30+(20+1*10))",
		},
		{
			name:         "1+1-(2.000000)",
			example:      "1+1-(2.000000)",
			expected_str: "(2.000000)",
		},
		{
			name:         "52",
			example:      "52",
			expected_str: "",
		},
		{
			name:         "operation withot value",
			example:      "+1",
			expected_str: "",
			expected_err: calc.OperationWithoutValue,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := calc.GetExampleNew(test.example)
			if !errors.Is(err, test.expected_err) {
				t.Errorf("GetExample(%s) got error %q, but expected %q", test.example, err, test.expected_err)
				return
			}

			if got.String != test.expected_str {
				t.Errorf("GetExample(%q) got %q, but expected %q", test.example, got.String, test.expected_str)
			}
		})
	}
}

func TestEraseExample(t *testing.T) {
	cases := []struct {
		name         string
		example      string
		expected_str string
	}{
		{
			name:         "1+1+1",
			example:      "1+1+1",
			expected_str: "2.000000+1",
		},
		{
			name:         "1+1-(1+1)",
			example:      "1+1-(1+1)",
			expected_str: "1+1-(2.000000)",
		},
		{
			name:         "1+1-(1+1)-1+1",
			example:      "1+1-(1+1)-1+1",
			expected_str: "1+1-(2.000000)-1+1",
		},
		{
			name:         "(1)",
			example:      "(1)",
			expected_str: "1.000000",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got_ex, err := calc.GetExampleNew(test.example)
			if err != nil {
				t.Errorf("GetExampleNew(%s) got %q", test.example, err)
			}

			local_ex := got_ex.String
			if strings.ContainsAny(got_ex.String, "()") {
				local_ex = got_ex.String[1 : len(got_ex.String)-1]
			}

			answ, err := calc.SolveExample(got_ex)
			if err != nil {
				t.Errorf("SolveExample(%s) got %q", local_ex, err)
				t.Log(got_ex.First_value, got_ex.Second_value, got_ex.Operation)
			}

			pri_idx := strings.Index(test.example, got_ex.String)
			got := calc.EraseExample(test.example, got_ex.String, pri_idx, answ)
			if got != test.expected_str {
				t.Errorf("EraseExample(%q, %q, %d, %f) = %q, but expected: %q", test.example, got_ex.String, pri_idx, answ, got, test.expected_str)
			}
		})
	}
}

func TestCalc(t *testing.T) {
	cases := []struct {
		name           string
		expression     string
		expected_value float64
		expected_err   error
	}{
		{
			name:           "simple addition",
			expression:     "1+1",
			expected_value: 2,
			expected_err:   nil,
		},
		{
			name:           "addition with negative value",
			expression:     "-3+1",
			expected_value: 0,
			expected_err:   calc.OperationWithoutValue,
		},
		{
			name:           "addition with 3 values",
			expression:     "1+1+1",
			expected_value: 3,
			expected_err:   nil,
		},
		{
			name:           "simple multiply",
			expression:     "1*1",
			expected_value: 1,
			expected_err:   nil,
		},
		{
			name:           "simple division",
			expression:     "1/1",
			expected_value: 1,
			expected_err:   nil,
		},
		{
			name:           "division with addition",
			expression:     "2+1*1",
			expected_value: 3,
			expected_err:   nil,
		},
		{
			name:           "hard example 1",
			expression:     "2+1*1+10/2",
			expected_value: 8,
			expected_err:   nil,
		},
		{
			name:           "brackets",
			expression:     "(1+1)/(1+1)",
			expected_value: 1,
			expected_err:   nil,
		},
		{
			name:           "hard example with brackets",
			expression:     "(1+10*(23-3)/2)-12",
			expected_value: 89,
			expected_err:   nil,
		},
		{
			name:           "unkown operator 1",
			expression:     "1&1",
			expected_value: 0,
			expected_err:   calc.ParseError,
		},
		{
			name:           "unkown operator 2",
			expression:     "1+&1",
			expected_value: 0,
			expected_err:   calc.ParseError,
		},
		{
			name:           "operation without value",
			expression:     "1+1*",
			expected_value: 0,
			expected_err:   calc.OperationWithoutValue,
		},
		{
			name:           "operation without value 2",
			expression:     "2+2**2",
			expected_value: 0,
			expected_err:   calc.ParseError,
		},
		{
			name:           "without closed bracket",
			expression:     "((2+2-*(2",
			expected_value: 0,
			expected_err:   calc.BracketsNotFound,
		},
		{
			name:           "nothing",
			expression:     "",
			expected_value: 0,
			expected_err:   calc.ExpressionEmpty,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := calc.Calc(test.expression)

			if !errors.Is(err, test.expected_err) {
				t.Errorf("Calc(%s) got %q, but expected %q", test.expression, err, test.expected_err)
				return
			}

			if got != test.expected_value {
				t.Errorf("Calc(%q) = %f, but expected - %f", test.expression, got, test.expected_value)
			}
		})
	}
}
