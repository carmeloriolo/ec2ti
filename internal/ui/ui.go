package ui

import (
	"log"

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

type Ui struct {
	Title    string
	Table    Table
	Handlers HandlerMap
	Screen   tcell.Screen
	cursor   int
	offset   int
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
