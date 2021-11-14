package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/carmeloriolo/ec2ti/internal/components"
)

func HandleShell(u *Ui) {
	table := u.Table.(*components.InstanceTable)
	u.Screen.Clear()
	err := u.Screen.Suspend()
	if err != nil {
		return
	}
	it := u.Table.(*components.InstanceTable).Instances[table.Cursor+table.Offset] // it
	port, user, pkey, err := startPrompt(it.Keyname)
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
		"-p",
		port,
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
