package ui

import "github.com/gdamore/tcell/v2"

const (
	KeySlash = 47
	KeyJ     = 106
	KeyK     = 107
)

func initKeys() {
	tcell.KeyNames[tcell.Key(KeySlash)] = "/"
	tcell.KeyNames[tcell.Key(KeySlash)] = "j"
	tcell.KeyNames[tcell.Key(KeySlash)] = "k"
}

func init() {
	initKeys()
}
