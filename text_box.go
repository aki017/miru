package miru

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

type TextBox struct {
	x    int
	y    int
	w    int
	h    int
	text string
}

func (tb *TextBox) SetPosition(x int, y int, w int, h int) {
	tb.x = x
	tb.y = y
	tb.w = w
	tb.h = h
}

type renderer struct {
	x     int
	y     int
	esc   int
	fg    int
	bg    int
	deco  int
	cache int
}

func (r *renderer) start() {
	r.esc = 1
	r.fg = int(termbox.ColorDefault)
	r.bg = int(termbox.ColorDefault)
}

func (r *renderer) set() {
	if r.cache < 10 {
		r.deco = r.cache
	} else if r.cache < 38 {
		r.fg = r.cache - 29
	} else {
		r.bg = r.cache - 39
	}
	r.cache = 0
}

func (r *renderer) consume(c rune) []rune {
	switch r.esc {
	case 1:
		if c == '[' {
			r.esc++
		} else {
			return []rune{'\x1b', c}
		}
	default:
		if c == ';' {
			r.set()
			r.esc++
		} else if c != 'm' {
			t, _ := strconv.Atoi(string(c))
			r.cache = r.cache*10 + t
		} else {
			r.set()
			r.cache = 0
			r.esc = 0
		}
	}
	return nil
}

func (r *renderer) newline() {
	r.x = 0
	r.y++
}

func (r *renderer) draw(x int, y int, w int, h int, c rune) {
	termbox.SetCell(x+r.x, y+r.y, c, termbox.Attribute(r.fg), termbox.Attribute(r.bg))
	r.x++
}

func (tb *TextBox) SetText(text string) {
	tb.text = text
}

func (tb TextBox) Draw() {
	var r renderer
	for _, c := range tb.text {
		if c == '\x1b' {
			r.start()
		} else if r.esc != 0 {
			r.consume(c)
		} else if c != '\n' {
			r.draw(tb.x, tb.y, tb.w, tb.h, c)
		} else {
			r.newline()
		}
	}
}

func pp(s []byte) (out string) {
	var r interface{}
	json.Unmarshal(s, &r)
	out = g(r, 0)
	return
}

func g(o interface{}, i int) (out string) {
	out = ""
	switch o.(type) {
	case string:
		out = strings.Repeat(" ", i*2) + o.(string) + "\n"
	case map[string]interface{}:
		r := o.(map[string]interface{})
		for k, v := range r {
			out += strings.Repeat(" ", i*2) + k + "\n"
			out += g(v, i+1)
		}
	}
	return
}
