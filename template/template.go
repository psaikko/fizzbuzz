package template

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

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

func FizzBuzz(from, to int) {
	for i := from; i < to; i += templateLines {
		fmt.Printf(templateString, i, i+1, i+3, i+6, i+7, i+10, i+12, i+13)
	}
}

func BufferedFizzBuzz(from, to int) {
	const bufSize = 65536
	bw := bufio.NewWriterSize(os.Stdout, bufSize)

	for i := from; i < to; i += templateLines {
		bytes := []byte(fmt.Sprintf(templateString, i, i+1, i+3, i+6, i+7, i+10, i+12, i+13))
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

		start := jobIndex*linesPerJob + 1
		end := (jobIndex + 1) * linesPerJob

		for i := start; i < end; i += templateLines {
			bytes := []byte(fmt.Sprintf(templateString, i, i+1, i+3, i+6, i+7, i+10, i+12, i+13))
			buffer = append(buffer, bytes...)
		}
		out <- buffer
		buffer, swapBuffer = swapBuffer, buffer
		buffer = buffer[:0]
	}

}

func ParallelFizzBuzz(_, to int) {
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

	for jobIndex*linesPerJob < to {
		for i := 0; i < nThreads; i++ {
			os.Stdout.Write(<-bufferChannels[i])
			jobChannels[i] <- jobIndex
			jobIndex++
		}
	}
}
