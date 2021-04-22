package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

func Render(u *Ui) {
	u.Screen.Sync()
	s := u.Screen
	s.Clear()
	renderTitleBox(u)
	renderTable(u)
	s.Show()
}

func renderTitleBox(u *Ui) {
	t := u.Title
	sw, sh := u.Screen.Size()
	u.DrawBox(1, 0, sw-1, sh-1)
	u.DrawStr(sw/2-len(t)/2-1, 0, tcell.StyleDefault, t)
}

func renderTable(u *Ui) {
	headers := u.Table.Headers()
	sw, sh := u.Screen.Size()
	n := len(headers)
	delta := sw / n
	w := 1
	for _, v := range headers {
		u.DrawStr(w, 1, styles[Header], v)
		w += delta
	}
	for i, v := range u.Table.(*InstanceTable).Instances {
		targetStyle := styles[v.State]
		w = 1
		if i+2 > sh-3 {
			return
		}
		if i+2 == u.cursor {
			targetStyle = styles[SelectedRow]
		}
		for _, s := range strings.Split(v.String(), " ") {
			u.DrawStr(w, i+2, targetStyle, s)
			// Fill gaps drawing blank chars
			for j := (w + len(s)); j < (w + delta); j++ {
				u.DrawStr(j, i+2, targetStyle, " ")
			}
			w += delta
		}
	}
}
