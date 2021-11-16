package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/carmeloriolo/ec2ti/internal/components"
	"github.com/gdamore/tcell/v2"
	"github.com/kyokomi/emoji"
)

var (
	searchIcon   = emoji.Sprintf(":detective::")
	searchPrefix = fmt.Sprintf("%s ", searchIcon)
	isValidKey   = regexp.MustCompile(`^[a-zA-Z0-9 +-=._:/@]+$`).MatchString
)

func updateCursor(t *components.InstanceTable, s string) {
	cursor, offset := t.GetPositionByInstanceName(s)
	t.SetCursor(cursor)
	t.SetOffset(offset)
}

func HandleSearch(u *Ui, e tcell.EventKey) {
	table := u.Table.(*components.InstanceTable)
	k := tcell.Key(e.Rune())
	switch k {
	case tcell.KeyEsc, tcell.KeyEnter:
		if u.searchMode {
			table.SetTitle(table.DefaultTitle(len(table.Instances)))
			if table.Cursor == -1 {
				table.Cursor++
			}
			u.searchMode = !u.searchMode
		}
	case KeySlash:
		if !u.searchMode {
			table.SetTitle(searchPrefix)
			u.searchMode = !u.searchMode
		} else {
			if len(table.Title) < 32 {
				title := fmt.Sprintf("%s%s", table.Title, string(k))
				table.SetTitle(title)
				updateCursor(table, strings.TrimPrefix(title, searchPrefix))
			}
		}
	case KeyBackspace:
		if len(table.Title) > len(searchPrefix) {
			title := table.Title[:len(table.Title)-1]
			table.SetTitle(title)
			updateCursor(table, strings.TrimPrefix(title, searchPrefix))
		}
	default:
		if isValidKey(fmt.Sprint(k)) {
			if len(table.Title) < 32 {
				title := fmt.Sprintf("%s%s", table.Title, string(k))
				table.SetTitle(title)
				updateCursor(table, strings.TrimPrefix(title, searchPrefix))
			}
		}
	}
	u.Render()
}
