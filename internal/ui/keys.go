package ui

import "github.com/gdamore/tcell/v2"

const (
	KeySlash = 47
	KeyD     = 100
	KeyS     = 115
	KeyJ     = 106
	KeyK     = 107
)

func initKeys() {
	tcell.KeyNames[tcell.Key(KeySlash)] = "/"
	tcell.KeyNames[tcell.Key(KeyJ)] = "j"
	tcell.KeyNames[tcell.Key(KeyK)] = "k"
	tcell.KeyNames[tcell.Key(KeyS)] = "s"
	tcell.KeyNames[tcell.Key(KeyD)] = "d"
}

func init() {
	initKeys()
}
