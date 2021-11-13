package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Handler func(*Ui)
type HandlerMap map[tcell.Key]Handler

func HandleCtrlC(u *Ui) {
	u.Screen.Fini()
	os.Exit(0)
}

func HandleSearch(u *Ui) {
	log.Println("TODO Search Instance")
	time.Sleep(time.Second * 2) // Give user some seconds to read the error
	u.Render()
}

func HandleDescribe(u *Ui) {
	log.Println("TODO Describe Instance")
	time.Sleep(time.Second * 2) // Give user some seconds to read the error
	u.Render()
}

func HandleShell(u *Ui) {
	table := u.Table.(*InstanceTable)
	u.Screen.Clear()
	err := u.Screen.Suspend()
	if err != nil {
		return
	}
	it := u.Table.(*InstanceTable).Instances[table.Cursor+table.Offset] // it
	user, pkey, err := startPrompt(it.Keyname)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %s", err.Error()))
		time.Sleep(time.Second * 2) // Give user some seconds to read the error
		err = u.Screen.Resume()
		if err == nil {
			u.Render()
		}
		return
	}
	host := fmt.Sprintf("%s@%s", user, it.Ip)
	pkey = fmt.Sprintf("~/.ssh/%s", pkey)
	cmd := exec.Command(
		"ssh",
		"-o ConnectTimeout=5",
		"-o StrictHostKeyChecking=no",
		"-i",
		pkey,
		"-t",
		host,
		"/bin/bash")
	// fmt.Println(cmd.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Println(fmt.Sprintf("Error: %s", err.Error()))
		time.Sleep(time.Second * 2) // Give user some seconds to read the error
	}
	err = cmd.Wait()
	if err != nil {
		log.Println(fmt.Sprintf("Error: %s", err.Error()))
		time.Sleep(time.Second * 2) // Give user some seconds to read the error
	}
	err = u.Screen.Resume()
	if err != nil {
		return
	}
	u.Render()
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
