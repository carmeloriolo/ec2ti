package ui

import (
	"log"
	"os"
	"time"

	"github.com/carmeloriolo/ec2ti/internal/components"
	"github.com/gdamore/tcell/v2"
)

type Handler func(*Ui, tcell.Key)
type HandlerMap map[tcell.Key]Handler

func HandleCtrlC(u *Ui, k tcell.Key) {
	u.Screen.Fini()
	os.Exit(0)
}

func HandleDescribe(u *Ui, k tcell.Key) {
	log.Println("TODO Describe Instance")
	table := u.Table.(*components.InstanceTable)
	it := u.Table.(*components.InstanceTable).Instances[table.Cursor+table.Offset] // it
	log.Println(it)
	time.Sleep(time.Second * 2) // Give user some seconds to read the error
	u.Render()
}

func HandleNavigateUp(u *Ui, k tcell.Key) {
	table := u.Table.(*components.InstanceTable)
	if table.Cursor > 0 {
		table.Cursor--
	} else if table.Offset > 0 {
		table.Offset--
	}
	u.Render()
}

func HandleNavigateDown(u *Ui, k tcell.Key) {
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
