package main

import (
	"fizzbuzz/baseline"
	"fizzbuzz/fixedwidth"
	"fizzbuzz/template"
	"fmt"
	"os"
)

const limit = 1 << 61

func main() {

	choice := "FixedWidth"
	if len(os.Args) > 1 {
		choice = os.Args[1]
	}

	strategies := map[string]func(int, int){
		"Baseline":           baseline.FizzBuzz,
		"Template":           template.FizzBuzz,
		"Buffered":           baseline.BufferedFizzBuzz,
		"BufferedTemplate":   template.BufferedFizzBuzz,
		"ParallelTemplate":   template.ParallelFizzBuzz,
		"FixedWidth":         fixedwidth.FizzBuzz,
		"ParallelFixedWidth": fixedwidth.ParallelFizzBuzz,
	}

	if f, ok := strategies[choice]; ok {
		f(1, limit)
	} else {
		fmt.Println("Strategy", choice, "not defined")
	}
}
