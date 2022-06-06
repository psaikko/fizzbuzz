package main

import (
	"fizzbuzz/baseline"
	"fizzbuzz/bufferedwriter"
	"fizzbuzz/fixedwidth"
	"fizzbuzz/template"
	"fmt"
	"os"
)

const limit = 1 << 61

func main() {

	choice := "ParallelFixedWidth"
	if len(os.Args) > 1 {
		choice = os.Args[1]
	}

	strategies := map[string]func(int, int){
		"Baseline":           baseline.FizzBuzz,           //  ~9 MiB/s
		"Template":           template.FizzBuzz,           //  80
		"BufferedWriter":     bufferedwriter.FizzBuzz,     // 204
		"BufferedTemplate":   template.BufferedFizzBuzz,   // 196
		"ParallelTemplate":   template.ParallelFizzBuzz,   // 420
		"FixedWidth":         fixedwidth.FizzBuzz,         // 350
		"ParallelFixedWidth": fixedwidth.ParallelFizzBuzz, // 800
	}

	if f, ok := strategies[choice]; ok {
		f(1, limit)
	} else {
		fmt.Println("Strategy", choice, "not defined")
	}
}
