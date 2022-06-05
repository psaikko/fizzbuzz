package bufferedwriter

import (
	"io"
	"os"
	"strconv"
)

type BufferedWriter struct {
	bufferSize int
	bufferPtr  int
	buffer     []byte
	f          io.Writer
}

func New(f io.Writer, bufferSize int) BufferedWriter {
	return BufferedWriter{
		bufferSize: bufferSize,
		bufferPtr:  0,
		buffer:     make([]byte, bufferSize),
		f:          f,
	}
}

func (bw *BufferedWriter) Write(data []byte) {
	if bw.bufferPtr+len(data) >= bw.bufferSize {
		bw.Flush()
	}

	copy(bw.buffer[bw.bufferPtr:], data)
	bw.bufferPtr += len(data)
}

func (bw *BufferedWriter) Flush() {
	bw.f.Write(bw.buffer[:bw.bufferPtr])
	bw.bufferPtr = 0
}

func FizzBuzz(from, to int) {
	fizzBuzzBytes := []byte("FizzBuzz\n")
	fizzBytes := []byte("Fizz\n")
	buzzBytes := []byte("Buzz\n")

	bw := New(os.Stdout, 65536)

	for i := from; i < to; i++ {
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
		bw.Write(bytes)
	}

	bw.Flush()
}
