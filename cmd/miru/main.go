package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/aki017/gq"
	"github.com/aki017/miru"
	"github.com/masatana/go-textdistance"
	"github.com/nsf/termbox-go"
)

type Like struct {
	s      *[]gq.Jv
	target string
}

func (l Like) Len() int {
	return len(*l.s)
}

func (l Like) Swap(i, j int) {
	(*l.s)[i], (*l.s)[j] = (*l.s)[j], (*l.s)[i]
}

func (l Like) Less(i, j int) bool {
	return textdistance.JaroWinklerDistance(gq.JvString((*l.s)[i]).StringValue(), l.target) > textdistance.JaroWinklerDistance(gq.JvString((*l.s)[j]).StringValue(), l.target)
}

type keyword string

func (k keyword) Like(s string) float64 {
	var v float64
	for i, c := range s {
		index := strings.Index(string(k), string(c))
		if index >= 0 {
			v += math.Pow(float64(1)/(math.Abs(float64(index-i))+1), 2) / math.Pow(float64(len(k)), 0.1)
		}
	}
	return v
}

func (k keyword) Preview(s string) string {
	a := ""
	for _, c := range k {
		if strings.Contains(string(s), string(c)) {
			a += "o"
		} else {
			a += "x"
		}
	}
	return a
}

func redraw_all() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	input.Draw()
	jq := gq.NewJQ()

	fp, _ := os.Open(os.Args[1])
	text, _ := ioutil.ReadAll(fp)
	result, _ := jq.Parse(string(text), string(input.Text))
	debug := "" //pp.Sprint(result)
	left.SetJv(result[0])
	left.Draw()

	//debug := ""

	if len(input.Text) > 0 {
		q := string(input.Text)
		b := q[strings.LastIndex(q, ".")+1:]
		q = q[:strings.LastIndex(q, ".")]
		q = strings.Replace(q, "[]", "[0]", -1)
		switch len(q) {
		case 0:
			q = "."
		default:
		}

		result, err := jq.Parse(string(text), q+" | keys")
		debug += fmt.Sprintf("%s\n %s %s\n", result, q, err)
		//debug += fmt.Sprintf("%s", err)
		var s []gq.Jv = gq.JvArray(result[0]).Array()
		sort.Sort(Like{s: &s, target: b})

		text := fmt.Sprintf("expr: %s, query: %s", q, b)

		//text += "["
		for _, k := range s {
			//text += string(k)
			debug += fmt.Sprintf("%3.3f %s\n", textdistance.JaroWinklerDistance(gq.JvString(k).StringValue(), b), gq.JvString(k).StringValue())
			//debug += k.Preview(b)
			//text += " "
		}
		//text += "]"
		suggest.SetText(text)
	} else {
		suggest.SetText("nothing")
	}
	right.SetText(debug)
	right.Draw()
	suggest.Draw()

	termbox.Flush()
}

var input miru.InputBox
var left miru.TreeView
var right miru.TextBox
var suggest miru.TextBox

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output256)
	termbox.SetInputMode(termbox.InputEsc)
	w, h := termbox.Size()
	for i := 0; i < h; i++ {
		const coldef = termbox.ColorDefault
		termbox.SetCell(w/2, i, '|', coldef, coldef)
	}
	input.SetPosition(0, 0, w/2, 1)
	input.SetPrompt("jq > ")
	input.SetText(".")
	suggest.SetPosition(w/2, 0, w/2, 1)
	left.SetPosition(0, 1, w/2-1, h)
	right.SetPosition(w/2, 1, w/2, h)

	redraw_all()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				//input.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				//input.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				input.DeleteRuneBackward()
			case termbox.KeyDelete, termbox.KeyCtrlD:
				//	input.DeleteRuneForward()
			case termbox.KeyTab:
				input.InsertRune('\t')
			case termbox.KeySpace:
				input.InsertRune(' ')
			case termbox.KeyCtrlK:
				////	input.DeleteTheRestOfTheLine()
			case termbox.KeyHome, termbox.KeyCtrlA:
				//	input.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				//	input.MoveCursorToEndOfTheLine()
			default:
				if ev.Ch != 0 {
					input.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw_all()
	}
}
