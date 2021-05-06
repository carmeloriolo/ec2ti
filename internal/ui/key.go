package ui

import "github.com/gdamore/tcell/v2"

const (
	KeySlash = 47
)

func initKeys() {
	tcell.KeyNames[tcell.Key(KeySlash)] = "/"
}

func init() {
	initKeys()
}
