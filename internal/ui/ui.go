package ui

import (
	"errors"
	"log"

	"github.com/carmeloriolo/ec2ti/internal/components"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

type ViewModeType int

const (
	ViewModeNormal   ViewModeType = 0
	ViewModeSearch   ViewModeType = 1
	ViewModeDescribe ViewModeType = 2
)

var (
	componentsRatio = 4
	defaultTitle    = "Title goes here"
	DefaultHandlers = HandlerMap{
		tcell.KeyCtrlC: HandleCtrlC,
		tcell.KeyUp:    HandleNavigateUp,
		tcell.KeyDown:  HandleNavigateDown,
		KeyEsc:         HandleEsc,
		KeyK:           HandleNavigateUp,
		KeyJ:           HandleNavigateDown,
		KeyD:           HandleDescribe,
		KeyS:           HandleShell,
		KeySlash:       HandleSearch,
	}
)

type Ui struct {
	Title    string
	Header   components.Header
	Table    components.Table
	Handlers HandlerMap
	Screen   tcell.Screen
	yTable   int
	ViewMode ViewModeType
}

func (u *Ui) Render() {
	s := u.Screen
	_, sh := u.Screen.Size()
	u.yTable = sh / componentsRatio
	u.Table.OnTableResize(u.NumberOfRowsDisplayed())
	// s.Clear()
	u.Table.Render(u.Screen, u.yTable)
	u.Header.Render(u.Title, u.Screen, u.yTable)
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
			// Search Mode
			if u.ViewMode == ViewModeSearch {
				if k == tcell.KeyCtrlC {
					(u.Handlers[tcell.KeyCtrlC])(u, *ev) // the program will exit
					continue
				}
				(u.Handlers[KeySlash])(u, *ev)
				continue
			}
			// Describe Mode
			if u.ViewMode == ViewModeDescribe {
				if k == KeyEsc {
					(u.Handlers[KeyEsc])(u, *ev)
				}
			}
			// Normal Mode
			if f, present := u.Handlers[k]; present {
				f(u, *ev)
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
		Title:    defaultTitle,
		ViewMode: ViewModeNormal,
		Screen:   s,
		yTable:   sh / componentsRatio,
	}
}
