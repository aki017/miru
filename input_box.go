package miru

import "github.com/nsf/termbox-go"

type InputBox struct {
	x      int
	y      int
	w      int
	h      int
	Text   []rune
	Prompt []rune
	cur    int
}

func (in *InputBox) InsertRune(r rune) {
	in.Text = append(in.Text, r)
}

func (in *InputBox) DeleteRuneBackward() {
	in.Text = in.Text[:len(in.Text)-1]
}

func (in *InputBox) SetPosition(x int, y int, w int, h int) {
	in.x = x
	in.y = y
	in.w = w
	in.h = h
}

func (in *InputBox) SetText(t string) {
	in.Text = []rune(t)
}

func (in *InputBox) SetPrompt(t string) {
	in.Prompt = []rune(t)
}

func (in InputBox) Draw() {
	for i, c := range append(in.Prompt, in.Text...) {
		const coldef = termbox.ColorDefault
		termbox.SetCell(in.x+i, in.y, c, coldef, coldef)
	}
}
