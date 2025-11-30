package paginate

import (
	"fmt"
	"testing"
)

func TestNewPaginator(t *testing.T) {
	page := New(testStrSlice)
	fmt.Printf("%#v\n", page.og)
}

var testStrSlice = []string{
	"My stepdad Derek married my dad",
	"when I was 9 years old. Now I'm",
	"13, so we've spent a decent amount",
	"of time together. He's a good guy.",
	"My dad isn't part of the picture,",
	"so it's been nice to have Derek",
}
