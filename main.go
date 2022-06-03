package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("./%s [STRATEGY]\n", os.Args[0])
	}

	strategies := map[string]func(){
		"baseline":  baseline,
		"withBufio": withBufio,
	}

	if f, ok := strategies[os.Args[1]]; ok {
		f()
	} else {
		fmt.Println("Strategy", os.Args[1], "not defined")
	}

}

// ~5.5 MiB/s
func baseline() {
	for i := 0; ; i++ {
		if (i%3 == 0) && (i%5 == 0) {
			fmt.Println("FizzBuzz")
		} else if i%3 == 0 {
			fmt.Println("Fizz")
		} else if i%5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Printf("%d\n", i)
		}
	}
}

// ~62 MiB/s
func withBufio() {
	b := bufio.NewWriter(os.Stdout)
	for i := 0; ; i++ {
		if (i%3 == 0) && (i%5 == 0) {
			b.WriteString("FizzBuzz\n")
		} else if i%3 == 0 {
			b.WriteString("Fizz\n")
		} else if i%5 == 0 {
			b.WriteString("Buzz\n")
		} else {
			b.WriteString(fmt.Sprintf("%d\n", i))
		}
	}
}
