package baseline

import (
	"bufio"
	"fmt"
	"os"
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

func BufferedFizzBuzz(from, to int) {
	bw := bufio.NewWriterSize(os.Stdout, 65536)
	for i := from; i < to; i++ {
		bw.WriteString(FizzBuzzLine(i))
	}
	bw.Flush()
}
