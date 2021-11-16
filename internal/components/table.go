package components

import "github.com/gdamore/tcell/v2"

type Table interface {
	Render(tcell.Screen, int)
	Columns() []string
	Rows() []string
	OnTableResize(int)
	SetTitle(string)
	SetCursor(int)
	SetOffset(int)
}
