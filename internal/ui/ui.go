package ui

import (
	"errors"
	"log"
	"strings"

	"github.com/carmeloriolo/ec2ti/internal/components"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/kyokomi/emoji"
)

const (
	HeaderRow     = "HeaderRow"
	TopRow        = "TopRow"
	Row           = "Row"
	SelectedRow   = "SelectedRow"
	StoppedRow    = "StoppedRow"
	CommandRow    = "Command"
	ctrlCLabel    = "<Ctrl+C> Exit"
	searchLabel   = "</> Search by Name"
	describeLabel = "<d> Describe"
	shellLabel    = "<s> Shell"
)

var (
	componentsRatio = 4
	defaultTitle    = "Title goes here"
	styles          = map[string]tcell.Style{
		"pending":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorGray),
		"running":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkCyan),
		"stopping":      tcell.StyleDefault.Italic(false).Foreground(tcell.ColorOrange),
		"stopped":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkGray),
		"shutting-down": tcell.StyleDefault.Italic(false).Foreground(tcell.ColorRed),
		"terminated":    tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkRed),
		HeaderRow:       tcell.StyleDefault.Bold(true).Foreground(tcell.ColorWhite),
		CommandRow:      tcell.StyleDefault.Bold(false).Foreground(tcell.ColorBlue),
		TopRow:          tcell.StyleDefault.Bold(true),
		SelectedRow:     tcell.StyleDefault.Bold(true).Foreground(tcell.ColorBlack.TrueColor()).Background(tcell.ColorDarkGray),
	}
	DefaultHandlers = HandlerMap{
		tcell.KeyCtrlC: HandleCtrlC,
		tcell.KeyUp:    HandleNavigateUp,
		tcell.KeyDown:  HandleNavigateDown,
		KeyK:           HandleNavigateUp,
		KeyJ:           HandleNavigateDown,
		KeyD:           HandleDescribe,
		KeyS:           HandleShell,
		KeySlash:       HandleSearch,
	}
	commandLabels = []string{
		searchLabel,
		describeLabel,
		shellLabel,
		ctrlCLabel,
	}
)

type Ui struct {
	Title      string
	Header     components.Header
	Table      components.Table
	Handlers   HandlerMap
	Screen     tcell.Screen
	yTable     int
	searchMode bool
}

func (u *Ui) Render() {
	s := u.Screen
	_, sh := u.Screen.Size()
	u.yTable = sh / componentsRatio
	u.Table.OnTableResize(u.NumberOfRowsDisplayed())
	s.Clear()
	renderTable(u)
	renderHeader(u)
	s.Sync()
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
		yTable:   u.yTable,
	}
	return u
}

func (u *Ui) SetHeader(h components.Header) *Ui {
	u = &Ui{
		Title:    u.Title,
		Header:   h,
		Table:    u.Table,
		Handlers: u.Handlers,
		Screen:   u.Screen,
		yTable:   u.yTable,
	}
	return u
}

func (u *Ui) SetTable(t components.Table) *Ui {
	u = &Ui{
		Title:    u.Title,
		Header:   u.Header,
		Table:    t,
		Handlers: u.Handlers,
		Screen:   u.Screen,
		yTable:   u.yTable,
	}
	return u
}

func (u *Ui) SetHandlers(h HandlerMap) *Ui {
	u = &Ui{
		Title:    u.Title,
		Header:   u.Header,
		Table:    u.Table,
		Handlers: h,
		Screen:   u.Screen,
		yTable:   u.yTable,
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
	_, sh := s.Size()
	return &Ui{
		Title:      defaultTitle,
		searchMode: false,
		Screen:     s,
		yTable:     sh / componentsRatio,
	}
}

func (u *Ui) Run() error {
	for {
		switch ev := u.Screen.PollEvent().(type) {
		case *tcell.EventError:
			continue
		case *tcell.EventResize:
			u.Render()
		case *tcell.EventKey:
			k := ev.Key()
			if ev.Rune() != 0 {
				k = tcell.Key(ev.Rune())
			}
			if !u.searchMode {
				if f, present := u.Handlers[k]; present {
					f(u, k)
				}
			} else {
				if k != tcell.KeyCtrlC {
					(u.Handlers[KeySlash])(u, k)
				} else {
					if f, present := u.Handlers[k]; present {
						f(u, k)
					}
				}
			}
		default:
			return errors.New("unexpected input")
		}
	}
}

func (u *Ui) NumberOfRowsDisplayed() int {
	_, sh := u.Screen.Size()
	return sh - u.yTable - 5
}

func renderHeader(u *Ui) {
	screen := u.Screen
	sw, _ := screen.Size()
	title := formatTitle(u.Title)
	components.DrawHeaderBox(screen, 0, 0, sw-1, u.yTable)
	components.DrawStr(screen, sw/2-len(title)/2-1, 0, styles[HeaderRow], title)
	nRows := u.yTable - 3
	for i, r := range u.Header.Rows() {
		if i < nRows {
			components.DrawStr(screen, 2, 2+i, styles[HeaderRow], r)
		} else {
			break
		}
	}
	if sw > 50 {
		for i, l := range commandLabels {
			if i < nRows {
				components.DrawStr(screen, sw-2-len(l), i+2, styles[CommandRow], l)
			} else {
				return
			}
		}
	}
}

func renderTable(u *Ui) {
	table := u.Table.(*components.InstanceTable)
	columns := u.Table.Columns()
	screen := u.Screen
	sw, sh := screen.Size()
	n := len(columns)
	delta := sw / n
	for delta < 21 {
		n--
		if n == 0 {
			return
		}
		delta = sw / n
	}
	w := 2
	components.DrawTableBox(screen, 0, u.yTable-1, sw-1, sh-2)
	components.DrawStr(screen, sw/2-len(table.Title)/2-1, u.yTable-1, tcell.StyleDefault, table.Title)

	for _, v := range columns[0:n] {
		components.DrawStr(screen, w, u.yTable+1, styles[TopRow], v)
		w += delta
	}
	tableInstances := table.Instances[table.Offset:len(table.Instances)]

	for i, v := range tableInstances {
		w = 2
		if i+2+u.yTable > sh-4 {
			return
		}
		targetStyle := styles[v.State]
		if i == table.Cursor {
			targetStyle = styles[SelectedRow]
		}
		for c, str := range strings.Split(v.String(), " ") {
			if c != n {
				components.DrawStr(screen, w, u.yTable+i+2, targetStyle, str)
				// Fill gaps drawing blank chars
				for j := (w + len(str)); j < (w + delta); j++ {
					components.DrawStr(screen, j, u.yTable+i+2, targetStyle, " ")
				}
				w += delta
			} else {
				break
			}
		}
	}
}

func formatTitle(t string) string {
	return emoji.Sprintf(" :rocket: %s :beer:", t)
}
