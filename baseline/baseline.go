package baseline

import (
	"fmt"
	"strconv"
)

func FizzBuzzLine(i int) string {
	if (i%3 == 0) && (i%5 == 0) {
		return "FizzBuzz\n"
	} else if i%3 == 0 {
		return "Fizz\n"
	} else if i%5 == 0 {
		return "Buzz\n"
	} else {
		return strconv.Itoa(i) + "\n"
	}
}

func FizzBuzz(from, to int) {
	for i := from; i < to; i++ {
		fmt.Print(FizzBuzzLine(i))
	}
}
