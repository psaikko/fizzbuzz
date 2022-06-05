package main

import (
	"bufio"
	"fizzbuzz/baseline"
	"fizzbuzz/bufferedwriter"
	"fizzbuzz/fixedwidth"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

const limit = 1 << 61

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("./%s [STRATEGY]\n", os.Args[0])
	}

	strategies := map[string]func(int, int){
		"baseline":                baseline.FizzBuzz,
		"withTemplate":            withTemplate,
		"withBufio":               withBufio,
		"withTemplateBufio":       withTemplateBufio,
		"writeBytes":              writeBytes,
		"BufferedWriter":          bufferedwriter.FizzBuzz,
		"templateBufferedWriter":  templateBufferedWriter,
		"parallelTemplateBuffers": parallelTemplateBuffers,
		"FixedWidth":              fixedwidth.FizzBuzz,
	}

	if f, ok := strategies[os.Args[1]]; ok {
		f(1, limit)
	} else {
		fmt.Println("Strategy", os.Args[1], "not defined")
	}

}

const templateLines = 15
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
func withTemplate(from, to int) {
	for i := from; i < to; i += templateLines {
		fmt.Printf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14)
	}
}

// ~101 MiB/s
func withBufio(from, to int) {
	b := bufio.NewWriter(os.Stdout)
	for i := from; i < to; i++ {
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
func withTemplateBufio(from, to int) {
	b := bufio.NewWriter(os.Stdout)
	for i := from; i < to; i += templateLines {
		b.WriteString(fmt.Sprintf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14))
	}
}

// ~21 MiB/s
func writeBytes(from, to int) {
	fizzBuzzBytes := []byte("FizzBuzz\n")
	fizzBytes := []byte("Fizz\n")
	buzzBytes := []byte("Buzz\n")
	f := os.Stdout

	for i := from; i < to; i++ {
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

// ~208 MiB/s
func templateBufferedWriter(from, to int) {
	const bufSize = 65536
	bw := bufferedwriter.New(os.Stdout, bufSize)

	for i := from; i < to; i += templateLines {
		bytes := []byte(fmt.Sprintf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14))
		bw.Write(bytes)
	}

	bw.Flush()
}

const jobSize = 10000

func worker(out chan<- []byte, in <-chan int) {
	linesPerJob := jobSize * templateLines
	buffer := make([]byte, 0, linesPerJob)
	swapBuffer := make([]byte, 0, linesPerJob)

	for jobIndex := range in {

		start := jobIndex * linesPerJob
		end := (jobIndex + 1) * linesPerJob

		for i := start; i < end; i += templateLines {
			bytes := []byte(fmt.Sprintf(templateString, i+1, i+2, i+4, i+7, i+8, i+11, i+13, i+14))
			buffer = append(buffer, bytes...)
		}
		out <- buffer
		buffer, swapBuffer = swapBuffer, buffer
		buffer = buffer[:0]
	}

}

// ~1 GiB/s
func parallelTemplateBuffers(_, to int) {
	nThreads := runtime.NumCPU()
	bufferChannels := make([]chan []byte, nThreads)
	jobChannels := make([]chan int, nThreads)

	for i := 0; i < nThreads; i++ {
		bufferChannels[i] = make(chan []byte)
		jobChannels[i] = make(chan int)
	}

	jobIndex := 0
	linesPerJob := jobSize * templateLines

	for i := 0; i < nThreads; i++ {
		go worker(bufferChannels[i], jobChannels[i])
		jobChannels[i] <- jobIndex
		jobIndex++
	}

	for jobIndex*linesPerJob < limit {
		for i := 0; i < nThreads; i++ {
			os.Stdout.Write(<-bufferChannels[i])
			jobChannels[i] <- jobIndex
			jobIndex++
		}
	}
}
