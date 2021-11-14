package ui

import (
	"github.com/gdamore/tcell/v2"
)

type UiInterface interface {
	Render()
	Run() error
	GetScreen() tcell.Screen
}
