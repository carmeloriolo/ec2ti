package ui

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/kyokomi/emoji"
)

const (
	HeaderRow   = "HeaderRow"
	TopRow      = "TopRow"
	Row         = "Row"
	SelectedRow = "SelectedRow"
	StoppedRow  = "StoppedRow"
)

var (
	oldSw           = 0
	oldSh           = 0
	componentsRatio = 3
	defaultTitle    = "Title goes here"
	styles          = map[string]tcell.Style{
		"pending":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorGray),
		"running":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkCyan),
		"stopping":      tcell.StyleDefault.Italic(false).Foreground(tcell.ColorOrange),
		"stopped":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkGray),
		"shutting-down": tcell.StyleDefault.Italic(false).Foreground(tcell.ColorRed),
		"terminated":    tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkRed),
		HeaderRow:       tcell.StyleDefault.Bold(true).Foreground(tcell.ColorWhite),
		TopRow:          tcell.StyleDefault.Bold(true),
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
	Title          string
	Header         HeaderInterface
	Table          TableInterface
	Handlers       HandlerMap
	Screen         tcell.Screen
	cursor         int
	offset         int
	tableStartingY int
}

func (u *Ui) Render() {
	u.Screen.Sync()
	s := u.Screen
	sw, sh := u.Screen.Size()
	if sh < oldSh {
		u.cursor = 2
		u.offset = 0
	}
	s.Clear()
	oldSh = sh
	oldSw = sw
	u.tableStartingY = sh / componentsRatio
	renderHeader(u)
	renderTable(u)
	s.Show()
}
func (u *Ui) GetScreen() tcell.Screen {
	return u.Screen
}

func (u *Ui) SetTitle(s string) *Ui {
	u = &Ui{
		Title:          s,
		Table:          u.Table,
		Handlers:       u.Handlers,
		Screen:         u.Screen,
		cursor:         u.cursor,
		offset:         u.offset,
		tableStartingY: u.tableStartingY,
	}
	return u
}

func (u *Ui) SetHeader(h HeaderInterface) *Ui {
	u = &Ui{
		Title:          u.Title,
		Header:         h,
		Table:          u.Table,
		Handlers:       u.Handlers,
		Screen:         u.Screen,
		cursor:         u.cursor,
		offset:         u.offset,
		tableStartingY: u.tableStartingY,
	}
	return u
}

func (u *Ui) SetTable(t TableInterface) *Ui {
	u = &Ui{
		Title:          u.Title,
		Header:         u.Header,
		Table:          t,
		Handlers:       u.Handlers,
		Screen:         u.Screen,
		cursor:         u.cursor,
		offset:         u.offset,
		tableStartingY: u.tableStartingY,
	}
	return u
}

func (u *Ui) SetHandlers(h HandlerMap) *Ui {
	u = &Ui{
		Title:          u.Title,
		Header:         u.Header,
		Table:          u.Table,
		Handlers:       h,
		Screen:         u.Screen,
		cursor:         u.cursor,
		offset:         u.offset,
		tableStartingY: u.tableStartingY,
	}
	return u

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
	_, sh := s.Size()
	return &Ui{
		Title:          defaultTitle,
		Screen:         s,
		cursor:         2,
		offset:         0,
		tableStartingY: sh / componentsRatio,
	}
}

func renderTitleBox(u *Ui) {
	t := u.Title
	s := u.Screen
	sw, sh := s.Size()
	DrawBox(s, 1, 0, sw-1, sh-1)
	DrawStr(s, sw/2-len(t)/2-1, 0, tcell.StyleDefault, t)
}

func renderHeader(u *Ui) {
	screen := u.Screen
	sw, _ := screen.Size()
	title := formatTitle(u.Title)
	for i, r := range u.Header.Rows() {
		DrawStr(screen, 1, 2+i, styles[HeaderRow], r)
	}
	DrawLine(screen, 0, 0, sw)
	DrawStr(screen, sw/2-len(title)/2-1, 0, styles[HeaderRow], title)
}

func renderTable(u *Ui) {
	columns := u.Table.Columns()
	screen := u.Screen
	sw, sh := screen.Size()
	n := len(columns)
	delta := sw / n
	w := 1
	// Render box around table
	tableTitle := emoji.Sprintf(" :computer: EC2 Instances (%d) ", len(u.Table.Rows())+u.offset)
	DrawBox(screen, 0, u.tableStartingY-1, sw-1, sh-1)
	DrawStr(screen, sw/2-len(tableTitle)/2-1, u.tableStartingY-1, tcell.StyleDefault, tableTitle)
	for _, v := range columns {
		DrawStr(screen, w, u.tableStartingY+1, styles[TopRow], v)
		w += delta
	}
	for i, v := range u.Table.(*InstanceTable).Instances {
		targetStyle := styles[v.State]
		w = 1
		if i+2+u.tableStartingY > sh-3 {
			return
		}
		if i+2 == u.cursor {
			targetStyle = styles[SelectedRow]
		}
		for _, str := range strings.Split(v.String(), " ") {
			DrawStr(screen, w, u.tableStartingY+i+2, targetStyle, str)
			// Fill gaps drawing blank chars
			for j := (w + len(str)); j < (w + delta); j++ {
				DrawStr(screen, j, u.tableStartingY+i+2, targetStyle, " ")
			}
			w += delta
		}
	}
}

func formatTitle(t string) string {
	return emoji.Sprintf(" :rocket: %s :beer:", t)
}
