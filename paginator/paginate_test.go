package paginate

import (
	"fmt"
	"testing"
)

func TestNewPaginator(t *testing.T) {
	page := New(testIntSlice,
		WithPerPage(10),
	)
	page.Loop = true
	fmt.Printf("page %v: %#v\n", page.Page, page.Prev())
	page.PrevPage()
	fmt.Printf("page %v: %#v\n", page.Page, page.Current())
	page.NextPage()
	fmt.Printf("page %v: %#v\n", page.Page, page.Current())

}

var testIntSlice = []int{1, 2, 3, 4, 5, 6}

var testStrSlice = []string{
	"My stepdad Derek married my dad",
	"when I was 9 years old. Now I'm",
	"13, so we've spent a decent amount",
	"of time together. He's a good guy.",
	"My dad isn't part of the picture,",
	"so it's been nice to have Derek",
}
