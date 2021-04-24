package ui

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

const (
	Header      = "Header"
	Row         = "Row"
	SelectedRow = "SelectedRow"
	StoppedRow  = "StoppedRow"
)

var (
	defaultTitle = "Title goes here"
	styles       = map[string]tcell.Style{
		"pending":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorGray),
		"running":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkCyan),
		"stopping":      tcell.StyleDefault.Italic(false).Foreground(tcell.ColorOrange),
		"stopped":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkGray),
		"shutting-down": tcell.StyleDefault.Italic(false).Foreground(tcell.ColorRed),
		"terminated":    tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkRed),
		Header:          tcell.StyleDefault.Bold(true),
		SelectedRow:     tcell.StyleDefault.Bold(true).Foreground(tcell.ColorBlack.TrueColor()).Background(tcell.ColorDarkGray),
	}
	DefaultHandlers = HandlerMap{
		tcell.KeyCtrlC: HandleCtrlC,
		tcell.KeyEnter: HandleEnter,
		tcell.KeyUp:    HandleNavigateUp,
		tcell.KeyDown:  HandleNavigateDown,
	}
)

type UiInterface interface {
	Render()
	GetScreen() tcell.Screen
}

type Ui struct {
	Title    string
	Table    Table
	Handlers HandlerMap
	Screen   tcell.Screen
	cursor   int
	offset   int
}

func (u *Ui) Render() {
	u.Screen.Sync()
	s := u.Screen
	s.Clear()
	renderTitleBox(u)
	renderTable(u)
	s.Show()
}

func (u *Ui) GetScreen() tcell.Screen {
	return u.Screen
}

func (u *Ui) SetTitle(s string) *Ui {
	u = &Ui{
		Title:    s,
		Table:    u.Table,
		Handlers: u.Handlers,
		Screen:   u.Screen,
		cursor:   u.cursor,
		offset:   u.offset,
	}
	return u
}

func (u *Ui) SetTable(t Table) *Ui {
	u = &Ui{
		Title:    u.Title,
		Table:    t,
		Handlers: u.Handlers,
		Screen:   u.Screen,
		cursor:   u.cursor,
		offset:   u.offset,
	}
	return u
}

func (u *Ui) SetHandlers(h HandlerMap) *Ui {
	u = &Ui{
		Title:    u.Title,
		Table:    u.Table,
		Handlers: h,
		Screen:   u.Screen,
		cursor:   u.cursor,
		offset:   u.offset,
	}
	return u

}

func (u *Ui) ResetPosition() {
	u.cursor = 2
	u.offset = 0
}

func NewUi() *Ui {
	encoding.Register()
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}
	s.SetStyle(
		tcell.StyleDefault.
			Background(tcell.ColorBlack.TrueColor()).
			Foreground(tcell.ColorWhite))
	return &Ui{
		Title:  defaultTitle,
		Screen: s,
		cursor: 2,
		offset: 0,
	}
}

func renderTitleBox(u *Ui) {
	t := u.Title
	s := u.Screen
	sw, sh := s.Size()
	DrawBox(s, 1, 0, sw-1, sh-1)
	DrawStr(s, sw/2-len(t)/2-1, 0, tcell.StyleDefault, t)
}

func renderTable(u *Ui) {
	headers := u.Table.Headers()
	screen := u.Screen
	sw, sh := screen.Size()
	n := len(headers)
	delta := sw / n
	w := 1
	for _, v := range headers {
		DrawStr(screen, w, 1, styles[Header], v)
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
		for _, str := range strings.Split(v.String(), " ") {
			DrawStr(screen, w, i+2, targetStyle, str)
			// Fill gaps drawing blank chars
			for j := (w + len(str)); j < (w + delta); j++ {
				DrawStr(screen, j, i+2, targetStyle, " ")
			}
			w += delta
		}
	}
}
