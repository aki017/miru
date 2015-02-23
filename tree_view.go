package miru

import (
	"github.com/aki017/gq"
	"github.com/nsf/termbox-go"
)

const (
	COL_INVALID termbox.Attribute = 205
	COL_NULL                      = 155
	COL_FALSE                     = 75
	COL_TRUE                      = 75
	COL_NUMBER                    = 215
	COL_STRING                    = 85
	COL_ARRAY                     = 135
	COL_OBJECT                    = 135
)

type TreeView struct {
	x  int
	y  int
	w  int
	h  int
	jv gq.Jv
}

type rend struct {
	x                  int
	y                  int
	w                  int
	h                  int
	offset_x, offset_y int
	fg                 termbox.Attribute
	bg                 termbox.Attribute
}

func (tb *TreeView) SetPosition(x int, y int, w int, h int) {
	tb.x = x
	tb.y = y
	tb.w = w
	tb.h = h
}

func (tv *TreeView) SetJv(jv gq.Jv) {
	tv.jv = jv
}

func (tv TreeView) Draw() {
	r := rend{
		w:        tv.w,
		h:        tv.h,
		offset_x: tv.x,
		offset_y: tv.y,
	}
	tv.drawr(&r, 0, tv.jv)
}

func (r *rend) newCol() {
	r.x++
	if r.x > r.w {
		r.newLine()
		r.x -= r.w
	}
}

func (r *rend) newColi(i int) {
	r.x += i
	if r.x > r.w {
		r.newLine()
		r.x -= r.w
	}
}

func (r *rend) newLine() {
	r.x = 0
	r.y++
}
func (r *rend) putIndent(depth int) {
	m := []rune{' ', ' '}
	for i := 0; i < depth*2; i++ {
		r.newCol()
		termbox.SetCell(r.offset_x+r.x, r.offset_y+r.y, m[i%2], r.fg, termbox.Attribute(235+((i%4)/2)*2))
	}
}

func (r *rend) putString(s string) {
	for _, c := range s {
		r.newCol()
		termbox.SetCell(r.offset_x+r.x, r.offset_y+r.y, c, r.fg, r.bg)
	}
}

func (tv TreeView) drawr(r *rend, depth int, jv gq.Jv) {
	r.bg = termbox.Attribute(termbox.ColorDefault)
	switch jv.Kind() {
	case gq.KIND_INVALID:
		r.putIndent(depth)
		r.fg = COL_INVALID
		r.putString("invalid")
		r.newLine()
	case gq.KIND_NULL:
		r.putIndent(depth)
		r.fg = COL_NULL
		r.putString("null")
		r.newLine()
	case gq.KIND_TRUE:
		r.putIndent(depth)
		r.fg = COL_TRUE
		r.putString("true")
		r.newLine()
	case gq.KIND_FALSE:
		r.putIndent(depth)
		r.fg = COL_FALSE
		r.putString("false")
		r.newLine()
	case gq.KIND_NUMBER:
		r.putIndent(depth)
		r.fg = COL_NUMBER
		r.putString(jv.String())
		r.newLine()
	case gq.KIND_STRING:
		r.putIndent(depth)
		r.fg = COL_STRING
		r.putString(gq.JvString(jv).StringValue())
		r.newLine()
	case gq.KIND_ARRAY:
		if onelinable(jv) {
			tv.drawInline(r, jv)
		} else {
			r.putIndent(depth)
			r.fg = COL_ARRAY
			r.putString("[")
			r.newLine()
			for _, value := range gq.JvArray(jv).Array() {
				tv.drawr(r, depth+1, value)
			}
			r.putIndent(depth)
			r.fg = COL_ARRAY
			r.putString("]")
		}
		r.newLine()
	case gq.KIND_OBJECT:
		if onelinable(jv) {
			tv.drawInline(r, jv)
		} else {
			gq.JvObject(jv).ForEach(func(key gq.Jv, value gq.Jv) {

				r.putIndent(depth)
				r.fg = COL_OBJECT
				r.putString(gq.JvString(key).StringValue())
				r.fg = COL_OBJECT
				r.putString(": ")

				if onelinable(value) {
					tv.drawInline(r, value)
					r.newLine()
				} else {
					r.newLine()
					tv.drawr(r, depth+1, value)
				}
			})
		}
	}
}

func (tv TreeView) drawInline(r *rend, jv gq.Jv) {
	switch jv.Kind() {
	case gq.KIND_INVALID:
		r.fg = COL_INVALID
		r.putString("invalid")
	case gq.KIND_NULL:
		r.fg = COL_NULL
		r.putString("null")
	case gq.KIND_TRUE:
		r.fg = COL_TRUE
		r.putString("true")
	case gq.KIND_FALSE:
		r.fg = COL_FALSE
		r.putString("false")
	case gq.KIND_NUMBER:
		r.fg = COL_NUMBER
		r.putString(jv.String())
	case gq.KIND_STRING:
		r.fg = COL_STRING
		r.putString(gq.JvString(jv).StringValue())
	case gq.KIND_ARRAY:
		r.fg = COL_ARRAY
		r.putString("[ ")
		for i, value := range gq.JvArray(jv).Array() {
			if i != 0 {
				r.fg = COL_ARRAY
				r.putString(", ")
			}
			tv.drawInline(r, value)
		}
		r.fg = COL_ARRAY
		r.putString("]")
	case gq.KIND_OBJECT:
		r.fg = COL_OBJECT
		r.putString("{ ")
		gq.JvObject(jv).ForEach(func(key gq.Jv, value gq.Jv) {
			r.fg = COL_OBJECT
			r.putString(gq.JvString(key).StringValue())
			r.fg = COL_OBJECT
			r.putString(": ")
			tv.drawInline(r, value)
			r.putString(" ")
		})
		r.fg = COL_OBJECT
		r.putString("}")
	}
}

func isvalue(jv gq.Jv) bool {
	switch jv.Kind() {
	case gq.KIND_INVALID:
		fallthrough
	case gq.KIND_NULL:
		fallthrough
	case gq.KIND_TRUE:
		fallthrough
	case gq.KIND_FALSE:
		fallthrough
	case gq.KIND_NUMBER:
		fallthrough
	case gq.KIND_STRING:
		return true
	}
	return false
}
func onelinable(jv gq.Jv) bool {
	switch jv.Kind() {
	case gq.KIND_INVALID:
		fallthrough
	case gq.KIND_NULL:
		fallthrough
	case gq.KIND_TRUE:
		fallthrough
	case gq.KIND_FALSE:
		fallthrough
	case gq.KIND_NUMBER:
		fallthrough
	case gq.KIND_STRING:
		return true
	case gq.KIND_ARRAY:
		if gq.JvArray(jv).Length() > 5 {
			return false
		}
		for _, value := range gq.JvArray(jv).Array() {
			if !isvalue(value) {
				return false
			}
		}
		return true
	case gq.KIND_OBJECT:
		if gq.JvObject(jv).Length() > 1 {
			return false
		}
		valid := true
		gq.JvObject(jv).ForEach(func(key gq.Jv, value gq.Jv) {
			if !isvalue(value) {
				valid = false
			}
		})
		return true
	}
	return false
}
