package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
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
		"parallelTemplateBuffers":  parallelTemplateBuffers,
	}

	if f, ok := strategies[os.Args[1]]; ok {
		f()
	} else {
		fmt.Println("Strategy", os.Args[1], "not defined")
	}

}

// ~20 MiB/s
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
const templateString = `%d
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
FizzBuzz
`

// ~116 MiB/s
func withTemplate() {
	for i := 0; i < limit; i += templateSize {
		fmt.Printf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14)
	}
}

// ~101 MiB/s
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

// ~156 MiB/s
func withTemplateBufio() {
	b := bufio.NewWriter(os.Stdout)
	for i := 0; i < limit; i += templateSize {
		b.WriteString(fmt.Sprintf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14))
	}
}

// ~21 MiB/s
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

// bufsize = 64:    ~81 MiB/s
// bufsize = 128:   110
// bufsize = 256:   151
// bufsize = 512:   119
// bufsize = 1024:  118
// bufsize = 2048:  121
// bufsize = 4096:  155
// bufsize = 8192:  182
// bufsize = 16384: 200
// bufsize = 32768: 208
// bufsize = 65536: 215
func writeByteBuffers() {
	fizzBuzzBytes := []byte("FizzBuzz\n")
	fizzBytes := []byte("Fizz\n")
	buzzBytes := []byte("Buzz\n")

	f := os.Stdout

	const bufSize = 65536
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

// ~208 MiB/s
func writeTemplateByteBuffers() {
	f := os.Stdout

	const bufSize = 65536
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

const jobSize = 10000

func worker(out chan<- []byte, in <-chan int) {
	linesPerJob := jobSize * templateSize
	buffer := make([]byte, 0, linesPerJob)
	swapBuffer := make([]byte, 0, linesPerJob)

	for jobIndex := range in {

		start := jobIndex * linesPerJob
		end := (jobIndex + 1) * linesPerJob

		for i := start; i < end; i += templateSize {
			bytes := []byte(fmt.Sprintf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14))
			buffer = append(buffer, bytes...)
		}
		out <- buffer
		buffer, swapBuffer = swapBuffer, buffer
		buffer = buffer[:0]
	}

}

// ~1 GiB/s
func parallelTemplateBuffers() {
	nThreads := runtime.NumCPU()
	bufferChannels := make([]chan []byte, nThreads)
	jobChannels := make([]chan int, nThreads)

	for i := 0; i < nThreads; i++ {
		bufferChannels[i] = make(chan []byte)
		jobChannels[i] = make(chan int)
	}

	jobIndex := 0

	for i := 0; i < nThreads; i++ {
		go worker(bufferChannels[i], jobChannels[i])
		jobChannels[i] <- jobIndex
		jobIndex++
	}

	for {
		for i := 0; i < nThreads; i++ {
			os.Stdout.Write(<-bufferChannels[i])
			jobChannels[i] <- jobIndex
			jobIndex++
		}
	}
}
