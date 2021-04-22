package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func (u *Ui) DrawStr(x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		u.Screen.SetContent(x, y, c, comb, style)
		x += w
	}
}

func (u *Ui) DrawBox(x, y, w, h int) {
	for i := x; i < w; i++ {
		u.Screen.SetContent(i, 0, tcell.RuneHLine, []rune{}, tcell.StyleDefault)
		u.Screen.SetContent(i, h-1, tcell.RuneHLine, []rune{}, tcell.StyleDefault)
	}
	for i := y; i < h; i++ {
		if i == 0 {
			u.Screen.SetContent(0, i, tcell.RuneULCorner, []rune{}, tcell.StyleDefault)
			u.Screen.SetContent(w, i, tcell.RuneURCorner, []rune{}, tcell.StyleDefault)
		} else if i == h-1 {
			u.Screen.SetContent(0, i, tcell.RuneLLCorner, []rune{}, tcell.StyleDefault)
			u.Screen.SetContent(w, i, tcell.RuneLRCorner, []rune{}, tcell.StyleDefault)
		} else {
			u.Screen.SetContent(0, i, tcell.RuneVLine, []rune{}, tcell.StyleDefault)
			u.Screen.SetContent(w, i, tcell.RuneVLine, []rune{}, tcell.StyleDefault)
		}
	}
}
