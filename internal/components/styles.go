package components

import "github.com/gdamore/tcell/v2"

const (
	// Header
	ctrlCLabel    = "<Ctrl+C> Exit"
	searchLabel   = "</> Search by Name"
	describeLabel = "<d> Describe"
	shellLabel    = "<s> Shell"
	// InstanceTable
	HeaderRow   = "HeaderRow"
	TopRow      = "TopRow"
	Row         = "Row"
	SelectedRow = "SelectedRow"
	StoppedRow  = "StoppedRow"
	CommandRow  = "Command"
)

var (
	styles = map[string]tcell.Style{
		"pending":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorGray),
		"running":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkCyan),
		"stopping":      tcell.StyleDefault.Italic(false).Foreground(tcell.ColorOrange),
		"stopped":       tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkGray),
		"shutting-down": tcell.StyleDefault.Italic(false).Foreground(tcell.ColorRed),
		"terminated":    tcell.StyleDefault.Italic(false).Foreground(tcell.ColorDarkRed),
		HeaderRow:       tcell.StyleDefault.Bold(true).Foreground(tcell.ColorWhite),
		CommandRow:      tcell.StyleDefault.Bold(false).Foreground(tcell.ColorBlue),
		TopRow:          tcell.StyleDefault.Bold(true),
		SelectedRow:     tcell.StyleDefault.Bold(true).Foreground(tcell.ColorBlack.TrueColor()).Background(tcell.ColorDarkGray),
	}
)
