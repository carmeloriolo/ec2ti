package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Handler func(*Ui)
type HandlerMap map[tcell.Key]Handler

func HandleCtrlC(u *Ui) {
	u.Screen.Fini()
	log.Println("CTRLC Received...Exiting")
	os.Exit(0)
}

func HandleSearch(u *Ui) {
	u.Screen.Fini()
	log.Println("Searching")
	os.Exit(0)
}

func HandleEnter(u *Ui) {
	table := u.Table.(*InstanceTable)
	u.Screen.Clear()
	u.Screen.Fini()
	it := u.Table.(*InstanceTable)
	fmt.Printf("Selected: %s\n", it.Instances[table.Cursor+table.Offset])
	os.Exit(0)
}

func HandleNavigateUp(u *Ui) {
	table := u.Table.(*InstanceTable)
	if table.Cursor > 0 {
		table.Cursor--
	} else if table.Offset > 0 {
		table.Offset--
	}
	u.Render()
}

func HandleNavigateDown(u *Ui) {
	table := u.Table.(*InstanceTable)
	n := table.RowsDisplayed
	if table.Cursor < n-1 {
		table.Cursor++
	} else {
		if table.Cursor+table.Offset < len(table.Instances)-1 {
			table.Offset++
		}
	}
	u.Render()
}
