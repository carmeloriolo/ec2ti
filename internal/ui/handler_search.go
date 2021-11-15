package ui

import (
	"fmt"
	"regexp"

	"github.com/carmeloriolo/ec2ti/internal/components"
	"github.com/gdamore/tcell/v2"
	"github.com/kyokomi/emoji"
)

var (
	searchIcon   = emoji.Sprintf(":detective:")
	searchPrefix = fmt.Sprintf("%s ", searchIcon)
	isValidKey   = regexp.MustCompile(`^[a-zA-Z0-9 +-=._:/@]+$`).MatchString
)

func HandleSearch(u *Ui, k tcell.Key) {
	table := u.Table.(*components.InstanceTable)
	switch k {
	case tcell.KeyEsc, tcell.KeyEnter:
		if u.searchMode {
			table.SetTitle(table.DefaultTitle(len(table.Instances)))
			u.searchMode = !u.searchMode
		}
	case KeySlash:
		if !u.searchMode {
			table.SetTitle(searchPrefix)
			u.searchMode = !u.searchMode
		} else {
			if len(table.Title) < 32 {
				table.SetTitle(fmt.Sprintf("%s%s", table.Title, string(k)))
			}
		}
	case KeyBackspace:
		if len(table.Title) > len(searchPrefix) {
			table.SetTitle(table.Title[:len(table.Title)-1])
		}
	default:
		if isValidKey(string(k)) {
			if len(table.Title) < 32 {
				table.SetTitle(fmt.Sprintf("%s%s", table.Title, string(k)))
			}
		}
	}
	table.SetCursor(3) // TODO Implement dynamic search for cursor and offset
	u.Render()
}
