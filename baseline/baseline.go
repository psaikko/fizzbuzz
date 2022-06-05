package baseline

import (
	"fmt"
	"strconv"
)

func FizzBuzzLine(i int) string {
	if (i%3 == 0) && (i%5 == 0) {
		return "FizzBuzz"
	} else if i%3 == 0 {
		return "Fizz"
	} else if i%5 == 0 {
		return "Buzz"
	} else {
		return strconv.Itoa(i) + "\n"
	}
}

func FizzBuzz(from, to int) {
	for i := from; i < to; i++ {
		fmt.Println(FizzBuzzLine(i))
	}
}
