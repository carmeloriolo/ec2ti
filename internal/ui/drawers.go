package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func DrawStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	sw, _ := s.Size()
	for _, c := range str {
		if x == sw-1 {
			return
		}
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			// comb = []rune{c}
			// c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func DrawHeaderBox(s tcell.Screen, x, y, w, h int) {
	for i := x; i < w; i++ {
		s.SetContent(i, y, tcell.RuneHLine, []rune{}, tcell.StyleDefault)
	}
	for i := y; i < h; i++ {
		if i == y {
			s.SetContent(x, i, tcell.RuneULCorner, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneURCorner, []rune{}, tcell.StyleDefault)
		} else if i == h {
			s.SetContent(x, i, tcell.RuneLLCorner, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneLRCorner, []rune{}, tcell.StyleDefault)
		} else {
			s.SetContent(x, i, tcell.RuneVLine, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneVLine, []rune{}, tcell.StyleDefault)
		}
	}
}

func DrawTableBox(s tcell.Screen, x, y, w, h int) {
	for i := x; i < w; i++ {
		s.SetContent(i, y, tcell.RuneHLine, []rune{}, tcell.StyleDefault)
		s.SetContent(i, h-1, tcell.RuneHLine, []rune{}, tcell.StyleDefault)
	}
	for i := y; i < h; i++ {
		if i == y {
			s.SetContent(x, i, tcell.RuneULCorner, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneURCorner, []rune{}, tcell.StyleDefault)
		} else if i == h-1 {
			s.SetContent(x, i, tcell.RuneLLCorner, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneLRCorner, []rune{}, tcell.StyleDefault)
		} else {
			s.SetContent(x, i, tcell.RuneVLine, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneVLine, []rune{}, tcell.StyleDefault)
		}
	}
}

func DrawLine(s tcell.Screen, x, y, l int) {
	for i := x; i < l; i++ {
		s.SetContent(i, y, tcell.RuneHLine, []rune{}, tcell.StyleDefault)
	}
}
