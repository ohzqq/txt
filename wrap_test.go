package txt

import (
	"fmt"
	"testing"
)

func TestWrapToSlice(t *testing.T) {
	lines := WrapToSlice(tstStr, 40)
	if len(lines) != 6 {
		fmt.Printf("%#v\n", len(lines))
	}
}

func TestWrapToSliceChunk(t *testing.T) {
	lines := WrapAndChunk(tstStr, 40, 3)
	if len(lines) != 2 {
		fmt.Printf("%#v\n", len(lines))
	}
}

const tstStr = `My stepdad Derek married my dad when I was 9 years old. Now I'm 13, so we've spent a decent amount of time together. He's a good guy. My dad isn't part of the picture, so it's been nice to have Derek around.`
