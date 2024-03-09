package sep

import (
	"fmt"
	"testing"
)

func TestSplit(t *testing.T) {
	var testStrings = map[string]int{
		`the quick, brown fox jumped! (and is running)`: 8,
	}
	for test, want := range testStrings {
		split := Split(test)
		fmt.Printf("got %d strings, wanted %d\n", len(split), want)
		for _, s := range split {
			fmt.Printf("split %v\n", s)
		}

		split = Split(test, Comma, Space)
		fmt.Printf("got %d strings, wanted %d\n", len(split), want)
		for _, s := range split {
			fmt.Printf("split %v\n", s)
		}

	}
}
