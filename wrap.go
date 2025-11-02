package txt

import (
	"strings"

	"github.com/duke-git/lancet/v2/slice"
)

const elipsis = `...`

func WrapToString(text string, w int) string {
	return strings.Join(WrapToSlice(text, w), "\n")
}

func WrapToSlice(str string, w int) []string {
	lines := split(str, w)
	return lines
}

func WrapAndChunk(str string, w, lh int) [][]string {
	return slice.Chunk(split(str, w), lh)
}

// borrowed from https://gist.github.com/AmrSaber/2468f546fb67dc31576a14e1209870e6
func split(str string, size int) []string {
	if size < 1 {
		return []string{}
	}

	str = strings.TrimSpace(str)

	start := 0
	chunks := make([]string, 0, len(str)/size)

	for start < len(str) {
		end := start + size
		if end >= len(str) {
			chunks = append(chunks, str[start:])
			break
		}

		// If the next character is a delimiter, take it
		// this is to avoid adding "-" at the end when the next character is a delimiter anyway
		next := str[end]
		if next == ' ' || next == '\n' {
			end++
		}

		chunk := str[start:end]
		cutWord := false

		// Try to find a new line within the limit
		length := strings.LastIndex(chunk, "\n")

		// If no new line found, try to find a space
		if length == -1 {
			length = strings.LastIndex(chunk, " ")

			// If no space found, then just split the text
			if length == -1 {
				length = len(chunk) - 1 // leave space for "-" character that will be appended
				cutWord = true
			}
		}

		chunk = chunk[:length]
		start += length
		if cutWord {
			chunk += "-"
		} else {
			// Ignore the space that we stopped at
			start++
		}

		chunks = append(chunks, chunk)
	}

	return chunks
}
