package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func DrawStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func DrawBox(s tcell.Screen, x, y, w, h int) {
	for i := x; i < w; i++ {
		s.SetContent(i, 0, tcell.RuneHLine, []rune{}, tcell.StyleDefault)
		s.SetContent(i, h-1, tcell.RuneHLine, []rune{}, tcell.StyleDefault)
	}
	for i := y; i < h; i++ {
		if i == 0 {
			s.SetContent(0, i, tcell.RuneULCorner, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneURCorner, []rune{}, tcell.StyleDefault)
		} else if i == h-1 {
			s.SetContent(0, i, tcell.RuneLLCorner, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneLRCorner, []rune{}, tcell.StyleDefault)
		} else {
			s.SetContent(0, i, tcell.RuneVLine, []rune{}, tcell.StyleDefault)
			s.SetContent(w, i, tcell.RuneVLine, []rune{}, tcell.StyleDefault)
		}
	}
}
