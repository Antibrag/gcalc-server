package main

import (
	"fmt"

	"github.com/Antibrag/gcalc-server/pkg/calc"
)

func main() {
	got, err := calc.Calc("1+1")
	fmt.Println(got, err)
}
