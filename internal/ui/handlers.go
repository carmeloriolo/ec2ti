package ui

import (
	"os"

	"github.com/carmeloriolo/ec2ti/internal/components"
	"github.com/gdamore/tcell/v2"
)

type Handler func(*Ui, tcell.EventKey)
type HandlerMap map[tcell.Key]Handler

func HandleEsc(u *Ui, e tcell.EventKey) {
	if u.ViewMode != ViewModeNormal {
		u.ViewMode = ViewModeNormal
		u.Screen.Clear()
		u.Render()
	}
}

func HandleCtrlC(u *Ui, e tcell.EventKey) {
	u.Screen.Fini()
	os.Exit(0)
}

func HandleDescribe(u *Ui, e tcell.EventKey) {
	// table := u.Table.(*components.InstanceTable)
	// it := u.Table.(*components.InstanceTable).Instances[table.Cursor+table.Offset] // it
	u.ViewMode = ViewModeDescribe
	u.Screen.Clear()
	x, y := u.Screen.Size()
	components.DrawBox(u.Screen, 0, 0, x-1, y-2)
	u.Screen.Sync()
}

func HandleNavigateUp(u *Ui, e tcell.EventKey) {
	table := u.Table.(*components.InstanceTable)
	if table.Cursor > 0 {
		table.Cursor--
	} else if table.Offset > 0 {
		table.Offset--
	}
	u.Render()
}

func HandleNavigateDown(u *Ui, e tcell.EventKey) {
	table := u.Table.(*components.InstanceTable)
	n := table.RowsDisplayed
	if table.Cursor < n-1 && table.Cursor < len(table.Instances)-1 {
		table.Cursor++
	} else {
		if table.Cursor+table.Offset < len(table.Instances)-1 {
			table.Offset++
		}
	}
	u.Render()
}
