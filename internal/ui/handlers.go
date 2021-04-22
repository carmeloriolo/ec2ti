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

func HandleEnter(u *Ui) {
	index := u.cursor - 2 + u.offset
	u.Screen.Clear()
	u.Screen.Fini()
	it := u.Table.(*InstanceTable)
	fmt.Printf("Selected: %s\n", it.Instances[index])
	os.Exit(0)
}

func HandleNavigateUp(u *Ui) {
	if u.cursor > 2 {
		u.cursor--
	} else {
		if u.offset > 0 {
			u.offset--
		}
	}
	it := u.Table.(*InstanceTable)
	Render(u.SetTable(&InstanceTable{
		Instances: it.Instances[u.offset:len(it.Instances)],
	}))
}

func HandleNavigateDown(u *Ui) {
	it := u.Table.(*InstanceTable)
	_, sh := u.Screen.Size()
	// Minus 3 because I have to consider top/bottom line + header line
	if u.cursor <= len(it.Instances) {
		if u.cursor < sh-3 {
			u.cursor++
		} else {
			if u.cursor+u.offset <= len(it.Instances) {
				u.offset++
			}
		}
		Render(u.SetTable(&InstanceTable{
			Instances: it.Instances[u.offset:len(it.Instances)],
		}))
	}
}
