package txt

import (
	"testing"
)

func TestDefaultAnalyzer(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           1,
		`quick brown fox`:     3,
		`QUICK BROWN FOX`:     3,
		`the quick brown fox`: 4,
		`the quick brown fox jumped and is running`:     8,
		`the quick, brown fox jumped! (and is running)`: 8,
		`The quick, brown fox jumped! (And is running)`: 8,
	}

	ana := NewAnalyzer()
	for test, want := range testStrings {
		tokens, err := ana.Tokenize(test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			t.Errorf("got %d tokens, wanted %d\n", numToks, want)
		}
	}
}

func TestAnalyzerToLower(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           1,
		`quick brown fox`:     3,
		`QUICK BROWN FOX`:     3,
		`the quick brown fox`: 4,
		`the quick brown fox jumped and is running`:     8,
		`the quick, brown fox jumped! (and is running)`: 8,
		`The quick, brown fox jumped! (And is running)`: 8,
	}

	ana := NewAnalyzer(
		ToLower,
	)
	for test, want := range testStrings {
		tokens, err := ana.Tokenize(test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			t.Errorf("got %d tokens, wanted %d\n", numToks, want)
		}
	}
}

func TestAnalyzerStripPunct(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           0,
		`quick brown fox`:     3,
		`QUICK BROWN FOX`:     3,
		`the quick brown fox`: 4,
		`the quick brown fox jumped and is running`:     8,
		`the quick, brown fox jumped! (and is running)`: 8,
		`The quick, brown fox jumped! (And is running)`: 8,
	}

	ana := NewAnalyzer(
		NoPunct,
	)
	for test, want := range testStrings {
		tokens, err := ana.Tokenize(test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			t.Errorf("%s: got %d tokens, wanted %d\n", test, numToks, want)
		}
	}
}
