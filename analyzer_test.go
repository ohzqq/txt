package txt

import (
	"strings"
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
		strings.ToLower,
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
		RmPunct,
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

func TestAnalyzerNormalize(t *testing.T) {
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
		strings.ToLower,
		RmPunct,
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

func TestAnalyzerStem(t *testing.T) {
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
		Stem,
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

func TestAnalyzerStopWords(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           1,
		`quick brown fox`:     3,
		`QUICK BROWN FOX`:     3,
		`the quick brown fox`: 3,
		`the quick brown fox jumped and is running`:     5,
		`the quick, brown fox jumped! (and is running)`: 6,
		`The quick, brown fox jumped! (And is running)`: 7,
	}

	ana := NewAnalyzer()
	ana.SetStopWords(DefaultStopWords())
	for test, want := range testStrings {
		tokens, err := ana.Tokenize(test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			println(test)
			for _, tok := range tokens {
				println(tok.Label)
			}
			t.Errorf("%s: got %d tokens, wanted %d\n", test, numToks, want)
		}
	}
}

func TestAnalyzerStopWordsNoPunct(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           0,
		`quick brown fox`:     3,
		`QUICK BROWN FOX`:     3,
		`the quick brown fox`: 3,
		`the quick brown fox jumped and is running`:     5,
		`the quick, brown fox jumped! (and is running)`: 5,
		`The quick, brown fox jumped! (And is running)`: 7,
	}

	ana := NewAnalyzer(
		RmPunct,
	)
	ana.SetStopWords(DefaultStopWords())
	for test, want := range testStrings {
		tokens, err := ana.Tokenize(test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			println(test)
			for _, tok := range tokens {
				println(tok.Label)
			}
			t.Errorf("%s: got %d tokens, wanted %d\n", test, numToks, want)
		}
	}
}

func TestAnalyzerStopWordsNormalize(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           0,
		`quick brown fox`:     3,
		`QUICK BROWN FOX`:     3,
		`the quick brown fox`: 3,
		`the quick brown fox jumped and is running`:     5,
		`the quick, brown fox jumped! (and is running)`: 5,
		`The quick, brown fox jumped! (And is running)`: 5,
	}

	ana := NewAnalyzer(
		strings.ToLower,
		RmPunct,
	)
	ana.SetStopWords(DefaultStopWords())
	for test, want := range testStrings {
		tokens, err := ana.Tokenize(test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			println(test)
			for _, tok := range tokens {
				println(tok.Label)
			}
			t.Errorf("%s: got %d tokens, wanted %d\n", test, numToks, want)
		}
	}
}
