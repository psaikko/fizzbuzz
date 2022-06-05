package fixedwidth

import (
	"fizzbuzz/baseline"
	"fizzbuzz/bufferedwriter"
	"fizzbuzz/ints"
	"fmt"
	"os"
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

		fmt.Println(rangeStart, rangeEnd)
		fmt.Println(templatesStart, templatesEnd)

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
		for i := ints.Max(templatesStart, templatesEnd+1); i < rangeEnd; i++ {
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
