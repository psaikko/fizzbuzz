package fixedwidth

import (
	"bufio"
	"fizzbuzz/baseline"
	"fizzbuzz/ints"
	"fmt"
	"io"
	"os"
	"runtime"
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

const cacheSize = 10000

var logCacheSize = ints.Log10(cacheSize)
var itoaCache = make([]string, cacheSize)

func initItoaCache() {
	// precompute string representations
	fmtString := fmt.Sprintf("%%0%dd", logCacheSize)
	for j := 0; j < cacheSize; j++ {
		itoaCache[j] = fmt.Sprintf(fmtString, j)
	}
}

type widthRange struct{ from, to, width int }

// getWidthRanges splits integer range [from,to] into disjoint ranges grouped by base10 representation length
func getWidthRanges(from, to int) []widthRange {
	ranges := []widthRange{}

	fromWidth := ints.Log10(from + 1)
	toWidth := ints.Log10(to + 1)

	for fromWidth < toWidth {
		toVal := ints.Pow(10, fromWidth) - 1
		ranges = append(ranges, widthRange{from, toVal, fromWidth})
		from = toVal + 1
		fromWidth += 1
	}

	ranges = append(ranges, widthRange{from, to, fromWidth})
	return ranges
}

func FizzBuzz(from, to int) {
	initItoaCache()

	const bufferSize = 65536
	bw := bufio.NewWriterSize(os.Stdout, bufferSize)

	for _, wr := range getWidthRanges(from, to) {
		// range which can be filled with templates
		templatesStart := ints.Min(wr.to, ints.CeilDiv(wr.from, templateLines)*templateLines)
		templatesEnd := templatesStart + ints.FloorDiv(wr.to-templatesStart+1, templateLines)*templateLines

		// handle values before first template
		for i := wr.from; i <= templatesStart; i++ {
			bw.Write([]byte(baseline.FizzBuzzLine(i)))
		}

		// print templates
		template, idxs := fixedWidthTemplate(wr.width)
		nextFlush := templatesStart

		for i := templatesStart; i < templatesEnd; i += templateLines {

			if i+14 > nextFlush || logCacheSize >= wr.width {
				// every $logCacheSize lines, write the entire buffer
				ints.CopyItoa(template, idxs[0]+wr.width, uint64(i+1))
				ints.CopyItoa(template, idxs[1]+wr.width, uint64(i+2))
				ints.CopyItoa(template, idxs[2]+wr.width, uint64(i+4))
				ints.CopyItoa(template, idxs[3]+wr.width, uint64(i+7))
				ints.CopyItoa(template, idxs[4]+wr.width, uint64(i+8))
				ints.CopyItoa(template, idxs[5]+wr.width, uint64(i+11))
				ints.CopyItoa(template, idxs[6]+wr.width, uint64(i+13))
				ints.CopyItoa(template, idxs[7]+wr.width, uint64(i+14))

				nextFlush += cacheSize
			} else {
				// write only the last $logCacheSize digits, others unchanged
				copy(template[idxs[0]+wr.width-logCacheSize:], itoaCache[(i+1)%cacheSize])
				copy(template[idxs[1]+wr.width-logCacheSize:], itoaCache[(i+2)%cacheSize])
				copy(template[idxs[2]+wr.width-logCacheSize:], itoaCache[(i+4)%cacheSize])
				copy(template[idxs[3]+wr.width-logCacheSize:], itoaCache[(i+7)%cacheSize])
				copy(template[idxs[4]+wr.width-logCacheSize:], itoaCache[(i+8)%cacheSize])
				copy(template[idxs[5]+wr.width-logCacheSize:], itoaCache[(i+11)%cacheSize])
				copy(template[idxs[6]+wr.width-logCacheSize:], itoaCache[(i+13)%cacheSize])
				copy(template[idxs[7]+wr.width-logCacheSize:], itoaCache[(i+14)%cacheSize])
			}
			bw.Write(template)
		}

		// handle values after last template
		for i := templatesEnd + 1; i <= wr.to; i++ {
			bw.Write([]byte(baseline.FizzBuzzLine(i)))
		}
	}

	bw.Flush()
}

func worker(in <-chan int, out chan<- []byte, templatesPerJob int, template []byte, width int, idxs []int) {
	buffer := make([]byte, len(template)*templatesPerJob)
	buffer2 := make([]byte, len(template)*templatesPerJob)
	buffer3 := make([]byte, len(template)*templatesPerJob)

	for i := 0; i < templatesPerJob; i++ {
		copy(buffer[len(template)*i:], template)
	}
	copy(buffer2, buffer)
	copy(buffer3, buffer)

	for jobLine := range in {

		nextFlush := (jobLine / cacheSize) * cacheSize
		t := false

		for i := 0; i < templatesPerJob; i++ {
			off := i * len(template)
			if i*templateLines+jobLine+13 > nextFlush || t {
				ints.CopyItoa(buffer, off+idxs[0]+width, uint64(i*templateLines+jobLine))
				ints.CopyItoa(buffer, off+idxs[1]+width, uint64(i*templateLines+jobLine+1))
				ints.CopyItoa(buffer, off+idxs[2]+width, uint64(i*templateLines+jobLine+3))
				ints.CopyItoa(buffer, off+idxs[3]+width, uint64(i*templateLines+jobLine+6))
				ints.CopyItoa(buffer, off+idxs[4]+width, uint64(i*templateLines+jobLine+7))
				ints.CopyItoa(buffer, off+idxs[5]+width, uint64(i*templateLines+jobLine+10))
				ints.CopyItoa(buffer, off+idxs[6]+width, uint64(i*templateLines+jobLine+12))
				ints.CopyItoa(buffer, off+idxs[7]+width, uint64(i*templateLines+jobLine+13))
				t = !t
				if t {
					nextFlush += cacheSize
				}
			} else {
				copy(buffer[off:], buffer[off-len(template):off])
				copy(buffer[off+idxs[0]+width-logCacheSize:], itoaCache[(i*templateLines+jobLine)%cacheSize])
				copy(buffer[off+idxs[1]+width-logCacheSize:], itoaCache[(i*templateLines+jobLine+1)%cacheSize])
				copy(buffer[off+idxs[2]+width-logCacheSize:], itoaCache[(i*templateLines+jobLine+3)%cacheSize])
				copy(buffer[off+idxs[3]+width-logCacheSize:], itoaCache[(i*templateLines+jobLine+6)%cacheSize])
				copy(buffer[off+idxs[4]+width-logCacheSize:], itoaCache[(i*templateLines+jobLine+7)%cacheSize])
				copy(buffer[off+idxs[5]+width-logCacheSize:], itoaCache[(i*templateLines+jobLine+10)%cacheSize])
				copy(buffer[off+idxs[6]+width-logCacheSize:], itoaCache[(i*templateLines+jobLine+12)%cacheSize])
				copy(buffer[off+idxs[7]+width-logCacheSize:], itoaCache[(i*templateLines+jobLine+13)%cacheSize])
			}
		}

		out <- buffer
		buffer, buffer2, buffer3 = buffer2, buffer3, buffer
	}

	close(out)
}

func writeParallel(f io.Writer, firstLine, lastLine, nWorkers, templatesPerJob int, template []byte, width int, placeholderIdxs []int) {

	totalLines := lastLine - firstLine + 1
	workerLines := templateLines * templatesPerJob
	linesPerRound := nWorkers * workerLines
	if totalLines%linesPerRound != 0 {
		panic("uneven work allocation")
	}

	jobChannels := make([]chan int, nWorkers)
	resultChannels := make([]chan []byte, nWorkers)
	totalJobs := ints.CeilDiv(totalLines, workerLines)
	jobsPerWorker := ints.CeilDiv(totalLines, workerLines*nWorkers)

	for i := 0; i < nWorkers; i++ {
		jobChannels[i] = make(chan int, jobsPerWorker*2)
		resultChannels[i] = make(chan []byte, 1)
		go worker(jobChannels[i], resultChannels[i], templatesPerJob, template, width, placeholderIdxs)
	}

	// deal out jobs to workers
	for job := 0; job < totalJobs; job++ {
		jobLine := firstLine + job*workerLines
		jobChannels[job%nWorkers] <- jobLine
	}

	// read buffers from workers
	for job := 0; job < totalJobs; job++ {
		f.Write(<-resultChannels[job%nWorkers])
	}
}

func ParallelFizzBuzz(from, to int) {
	initItoaCache()

	for _, wr := range getWidthRanges(from, to) {
		// range which can be filled with templates
		templatesStart := ints.Min(wr.to, ints.CeilDiv(wr.from, templateLines)*templateLines)
		templatesEnd := templatesStart + ints.FloorDiv(wr.to-templatesStart+1, templateLines)*templateLines

		// handle values before first template
		for i := wr.from; i <= templatesStart; i++ {
			os.Stdout.WriteString(baseline.FizzBuzzLine(i))
		}

		// write large chunks in parallel
		const templatesPerJob = 250
		template, placeholderIdxs := fixedWidthTemplate(wr.width)
		nWorkers := runtime.NumCPU()
		chunkSize := nWorkers * templateLines * templatesPerJob

		chunksStart := templatesStart
		chunksEnd := chunksStart + ints.FloorDiv(templatesEnd-templatesStart+1, chunkSize)*chunkSize

		if chunksEnd > templatesStart {
			writeParallel(os.Stdout, chunksStart+1, chunksEnd, nWorkers, templatesPerJob, template, wr.width, placeholderIdxs)
		}

		// handle values after last chunk
		for i := chunksEnd + 1; i <= wr.to; i++ {
			os.Stdout.WriteString(baseline.FizzBuzzLine(i))
		}
	}
}
