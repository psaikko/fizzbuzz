package fixedwidth

import (
	"fizzbuzz/baseline"
	"fizzbuzz/bufferedwriter"
	"fizzbuzz/ints"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
)

const templateLines = 15

func fixedWidthTemplate(valueWidth int) ([]byte, []int) {
	template := make([]byte, 0, 15+valueWidth*8+4*8)
	formatString := fmt.Sprintf("%%0%dd\n", valueWidth)
	placeholder := []byte(fmt.Sprintf(formatString, 0))
	placeholderIdxs := make([]int, 0, 8)
	fizzBytes := []byte("Fizz\n")
	buzzBytes := []byte("Buzz\n")
	fizzBuzzBytes := []byte("FizzBuzz\n")

	placeholderIdxs = append(placeholderIdxs, len(template))
	template = append(template, placeholder...)
	placeholderIdxs = append(placeholderIdxs, len(template))
	template = append(template, placeholder...)
	template = append(template, fizzBytes...)
	placeholderIdxs = append(placeholderIdxs, len(template))
	template = append(template, placeholder...)
	template = append(template, buzzBytes...)
	template = append(template, fizzBytes...)
	placeholderIdxs = append(placeholderIdxs, len(template))
	template = append(template, placeholder...)
	placeholderIdxs = append(placeholderIdxs, len(template))
	template = append(template, placeholder...)
	template = append(template, fizzBytes...)
	template = append(template, buzzBytes...)
	placeholderIdxs = append(placeholderIdxs, len(template))
	template = append(template, placeholder...)
	template = append(template, fizzBytes...)
	placeholderIdxs = append(placeholderIdxs, len(template))
	template = append(template, placeholder...)
	placeholderIdxs = append(placeholderIdxs, len(template))
	template = append(template, placeholder...)
	template = append(template, fizzBuzzBytes...)
	return template, placeholderIdxs
}

func FizzBuzz(from, to int) {

	bw := bufferedwriter.New(os.Stdout, 65536)

	rangeStart := from
	rangeEnd := ints.Pow(10, ints.Log10(rangeStart)+1) - 1
	rangeEnd = ints.Min(rangeEnd, to)

	for width := 1; ; width++ {
		// range which can be filled with templates
		templatesStart := ((rangeStart + templateLines - 1) / templateLines) * templateLines
		templatesEnd := (rangeEnd / templateLines) * templateLines

		// handle values before first template
		for i := rangeStart; i <= ints.Min(templatesStart, rangeEnd); i++ {
			bw.Write([]byte(baseline.FizzBuzzLine(i)))
		}

		// print templates
		template, placeholderIdxs := fixedWidthTemplate(width)
		for i := templatesStart; i < templatesEnd; i += templateLines {
			copy(template[placeholderIdxs[0]:], strconv.Itoa(i+1))
			copy(template[placeholderIdxs[1]:], strconv.Itoa(i+2))
			copy(template[placeholderIdxs[2]:], strconv.Itoa(i+4))
			copy(template[placeholderIdxs[3]:], strconv.Itoa(i+7))
			copy(template[placeholderIdxs[4]:], strconv.Itoa(i+8))
			copy(template[placeholderIdxs[5]:], strconv.Itoa(i+11))
			copy(template[placeholderIdxs[6]:], strconv.Itoa(i+13))
			copy(template[placeholderIdxs[7]:], strconv.Itoa(i+14))
			bw.Write(template)
		}

		// handle values after last template
		for i := ints.Max(templatesStart, templatesEnd+1); i <= rangeEnd; i++ {
			bw.Write([]byte(baseline.FizzBuzzLine(i)))
		}

		// update ranges
		rangeStart *= 10
		rangeEnd = ints.Min(rangeStart*10-1, to)
		if rangeStart > rangeEnd {
			break
		}
	}

	bw.Flush()
}

func worker(in <-chan int, out chan<- []byte, templatesPerJob int, template []byte, placeholderIdxs []int) {
	buffer := make([]byte, len(template)*templatesPerJob)
	swapBuffer := make([]byte, len(template)*templatesPerJob)

	for i := 0; i < templatesPerJob; i++ {
		copy(buffer[len(template)*i:], template)
	}
	copy(swapBuffer, buffer)

	for jobLine := range in {
		for i := 0; i < templatesPerJob; i++ {
			posOffset := i * len(template)
			copy(buffer[posOffset+placeholderIdxs[0]:], strconv.Itoa(jobLine+15*i))
			copy(buffer[posOffset+placeholderIdxs[1]:], strconv.Itoa(jobLine+15*i+1))
			copy(buffer[posOffset+placeholderIdxs[2]:], strconv.Itoa(jobLine+15*i+3))
			copy(buffer[posOffset+placeholderIdxs[3]:], strconv.Itoa(jobLine+15*i+6))
			copy(buffer[posOffset+placeholderIdxs[4]:], strconv.Itoa(jobLine+15*i+7))
			copy(buffer[posOffset+placeholderIdxs[5]:], strconv.Itoa(jobLine+15*i+10))
			copy(buffer[posOffset+placeholderIdxs[6]:], strconv.Itoa(jobLine+15*i+12))
			copy(buffer[posOffset+placeholderIdxs[7]:], strconv.Itoa(jobLine+15*i+13))
		}

		out <- buffer
		buffer, swapBuffer = swapBuffer, buffer
	}

	close(out)
}

func writeParallel(f io.Writer, firstLine, lastLine, nWorkers, templatesPerJob int, template []byte, placeholderIdxs []int) {

	totalLines := lastLine - firstLine + 1
	workerLines := templateLines * templatesPerJob
	linesPerRound := nWorkers * workerLines
	if totalLines%linesPerRound != 0 {
		panic("uneven allocation")
	}

	jobChannels := make([]chan int, 0)
	resultChannels := make([]chan []byte, 0)

	workersPos := firstLine

	for i := 0; i < nWorkers; i++ {
		jobChan := make(chan int)
		resultChan := make(chan []byte)
		go worker(jobChan, resultChan, templatesPerJob, template, placeholderIdxs)
		jobChannels = append(jobChannels, jobChan)
		resultChannels = append(resultChannels, resultChan)
		jobChan <- workersPos
		workersPos += workerLines
	}

	// deal out jobs to workers
	ctr := 0
	for ; workersPos < lastLine; workersPos += workerLines {
		f.Write(<-resultChannels[ctr%nWorkers])
		jobChannels[ctr%nWorkers] <- workersPos
		ctr++
	}

	// take last buffers and close channels
	for i := 0; i < nWorkers; i++ {
		f.Write(<-resultChannels[i])
		close(jobChannels[i])
	}
}

func ParallelFizzBuzz(from, to int) {

	f := os.Stdout

	rangeStart := from
	rangeEnd := ints.Pow(10, ints.Log10(rangeStart)+1) - 1
	rangeEnd = ints.Min(rangeEnd, to)

	for width := 1; ; width++ {
		// range which can be filled with templates
		templatesStart := ((rangeStart+templateLines-1)/templateLines)*templateLines + 1
		templatesEnd := (rangeEnd / templateLines) * templateLines
		nTemplatedLines := templatesEnd - templatesStart + 1

		// handle values before first template
		for i := rangeStart; i < ints.Min(templatesStart, rangeEnd+1); i++ {
			f.WriteString(baseline.FizzBuzzLine(i))
		}

		// write large chunks in parallel
		const templatesPerJob = 10000
		template, placeholderIdxs := fixedWidthTemplate(width)
		nWorkers := runtime.NumCPU()
		chunkSize := nWorkers * templateLines * templatesPerJob
		chunksEnd := templatesStart + (nTemplatedLines/chunkSize)*chunkSize - 1

		if chunksEnd > templatesStart {
			writeParallel(f, templatesStart, chunksEnd, nWorkers, templatesPerJob, template, placeholderIdxs)
		}

		// handle values after last chunk
		for i := chunksEnd + 1; i <= rangeEnd; i++ {
			f.WriteString(baseline.FizzBuzzLine(i))
		}

		// update ranges
		rangeStart *= 10
		rangeEnd = ints.Min(rangeStart*10-1, to)
		if rangeStart > rangeEnd {
			break
		}
	}
}
