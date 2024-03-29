package txt

import (
	"testing"

	"github.com/ohzqq/txt/sep"
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

	ana := New()
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

func TestKeywordAnalyzer(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           1,
		`quick brown fox`:     1,
		`QUICK BROWN FOX`:     1,
		`the quick brown fox`: 1,
		`the quick brown fox jumped and is running`:     1,
		`the quick, brown fox jumped! (and is running)`: 1,
		`The quick, brown fox jumped! (And is running)`: 1,
	}

	ana := New().Keywords()
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			t.Errorf("got %d tokens, wanted %d\n", numToks, want)
		}
	}
}

func TestFuzzySearch(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           1,
		`quick brown fox`:     1,
		`QUICK BROWN FOX`:     1,
		`the quick brown fox`: 1,
		`the quick brown fox jumped and is running`:     1,
		`the quick, brown fox jumped! (and is running)`: 1,
		`The quick, brown fox jumped! (And is running)`: 1,
	}

	ana := New().Keywords()
	//ana.Keywords()
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			t.Errorf("got %d tokens, wanted %d\n", numToks, want)
		}
		m := tokens.FuzzyFind("row")
		switch test {
		case ``, `!@@#$$%%^`:
			if m.Len() != 0 {
				t.Errorf("got %d tokens, wanted %d\n", m.Len(), want)
			}
		default:
			if want != m.Len() {
				t.Errorf("got %d tokens, wanted %d\n", m.Len(), want)
			}
		}
	}
}

func TestFuzzySearchNormalized(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           0,
		`quick brown fox`:     3,
		`QUICK BROWN FOX`:     3,
		`the quick brown fox`: 4,
		`the quick brown fox jumped and is running`:      8,
		`the quick, brown fox jumped! (and is running)`:  8,
		`The quick, brown fox jumping! (And is running)`: 8,
	}

	ana := Normalize()
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
		if err != nil {
			println(err.Error())
		}
		numToks := len(tokens)
		if want != numToks {
			t.Errorf("%s: got %d tokens, wanted %d\n", test, numToks, want)
		}

		m := tokens.FuzzyFind("ing")
		switch test {
		case `the quick brown fox jumped and is running`, `the quick, brown fox jumped! (and is running)`:
			if m.Len() != 1 {
				t.Errorf("got %d tokens, wanted %d\n", m.Len(), want)
			}
		case `The quick, brown fox jumping! (And is running)`:
			if m.Len() != 2 {
				t.Errorf("got %d tokens, wanted %d\n", m.Len(), want)
			}
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

	ana := New(ToLower)
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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

	ana := New(WithoutPunct)
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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

	ana := Normalize()
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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

	ana := New(WithStemmer)
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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
		`The quick, brown fox jumped! (And is running)`: 6,
	}

	ana := New(
		WithStopWords(DefaultStopWords()),
	)
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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
		`The quick, brown fox jumped! (And is running)`: 5,
	}

	ana := New(WithoutPunct).
		SetStopWords(DefaultStopWords())
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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

	ana := Normalize()
	ana.SetStopWords(DefaultStopWords())
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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

func TestAnalyzerKitchenSink(t *testing.T) {
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

	ana := Normalize(WithStemmer, WithDefaultStopWords)
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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

func TestAnalyzerSetFieldsFunc(t *testing.T) {
	var testStrings = map[string]int{
		``:                    0,
		`!@@#$$%%^`:           1,
		`quick brown fox`:     1,
		`QUICK BROWN FOX`:     1,
		`the quick brown fox`: 1,
		`the quick brown fox jumped and is running`:     1,
		`the quick, brown fox jumped! (and is running)`: 2,
		`The quick, brown fox jumped! (And is running)`: 2,
	}

	ana := New().WithSep(sep.Comma)
	for test, want := range testStrings {
		tokens, err := Tokenize(ana, test)
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

func TestUnicode(t *testing.T) {
	//for _, ra := range unicode.S.R16 {
	//fmt.Printf("%#U\n", int32(ra.Lo))
	//}
	//  for _, ra := range unicode.Punct.R32 {
	//    fmt.Printf("%#U\n", int32(ra.Lo))
	//  }

}
