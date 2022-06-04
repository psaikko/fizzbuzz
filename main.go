package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const limit = 1 << 32

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("./%s [STRATEGY]\n", os.Args[0])
	}

	strategies := map[string]func(){
		"baseline":                 baseline,
		"withTemplate":             withTemplate,
		"withBufio":                withBufio,
		"withTemplateBufio":        withTemplateBufio,
		"writeBytes":               writeBytes,
		"writeByteBuffers":         writeByteBuffers,
		"writeTemplateByteBuffers": writeTemplateByteBuffers,
	}

	if f, ok := strategies[os.Args[1]]; ok {
		f()
	} else {
		fmt.Println("Strategy", os.Args[1], "not defined")
	}

}

// ~5.5 MiB/s
func baseline() {
	for i := 0; i < limit; i++ {
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

const templateSize = 15
const templateString = `FizzBuzz
%d
%d
Fizz
%d
Buzz
Fizz
%d
%d
Fizz
Buzz
%d
Fizz
%d
%d
`

// ~45 MiB/s
func withTemplate() {
	for i := 0; i < limit; i += templateSize {
		fmt.Printf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14)
	}
}

// ~62 MiB/s
func withBufio() {
	b := bufio.NewWriter(os.Stdout)
	for i := 0; i < limit; i++ {
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
	b.Flush()
}

// ~104 MiB/s
func withTemplateBufio() {
	b := bufio.NewWriter(os.Stdout)
	for i := 0; i < limit; i += templateSize {
		b.WriteString(fmt.Sprintf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14))
	}
}

// ~5.9 MiB/s
func writeBytes() {
	fizzBuzzBytes := []byte("FizzBuzz\n")
	fizzBytes := []byte("Fizz\n")
	buzzBytes := []byte("Buzz\n")
	f := os.Stdout

	for i := 0; i < limit; i++ {
		if (i%3 == 0) && (i%5 == 0) {
			f.Write(fizzBuzzBytes)
		} else if i%3 == 0 {
			f.Write(fizzBytes)
		} else if i%5 == 0 {
			f.Write(buzzBytes)
		} else {
			f.Write([]byte(strconv.Itoa(i) + "\n"))
		}
	}
}

// bufsize = 64:    ~29 MiB/s
// bufsize = 128:    46
// bufsize = 256:    57
// bufsize = 512:    70
// bufsize = 1024:   88
// bufsize = 2048:  101
// bufsize = 4096:  114
// bufsize = 8192:  119
// bufsize = 16384: 119
func writeByteBuffers() {
	fizzBuzzBytes := []byte("FizzBuzz\n")
	fizzBytes := []byte("Fizz\n")
	buzzBytes := []byte("Buzz\n")

	f := os.Stdout

	const bufSize = 16384
	bufPtr := 0
	buffer := make([]byte, bufSize)

	for i := 0; i < limit; i++ {
		var bytes []byte
		if (i%3 == 0) && (i%5 == 0) {
			bytes = fizzBuzzBytes
		} else if i%3 == 0 {
			bytes = fizzBytes
		} else if i%5 == 0 {
			bytes = buzzBytes
		} else {
			bytes = []byte(strconv.Itoa(i) + "\n")
		}

		if bufPtr+len(bytes) >= bufSize {
			f.Write(buffer[:bufPtr])
			bufPtr = 0
		}

		copy(buffer[bufPtr:], bytes)
		bufPtr += len(bytes)
	}

	f.Write(buffer[:bufPtr])
}

// ~99 MiB/s
func writeTemplateByteBuffers() {
	f := os.Stdout

	const bufSize = 16384
	bufPtr := 0
	buffer := make([]byte, bufSize)

	for i := 0; i < limit; i += templateSize {
		bytes := []byte(fmt.Sprintf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14))

		if bufPtr+len(bytes) >= bufSize {
			f.Write(buffer[:bufPtr])
			bufPtr = 0
		}

		copy(buffer[bufPtr:], bytes)
		bufPtr += len(bytes)
	}

	f.Write(buffer[:bufPtr])
}
