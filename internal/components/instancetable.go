package components

import (
	"reflect"
	"strings"

	"github.com/carmeloriolo/ec2ti/internal/client"
	"github.com/gdamore/tcell/v2"
	"github.com/kyokomi/emoji"
)

type InstanceTable struct {
	Instances     []client.Instance
	Cursor        int
	Offset        int
	RowsDisplayed int
	Title         string
}

func (t *InstanceTable) Columns() []string {
	columns := []string{}
	v := reflect.ValueOf(client.Instance{})
	for i := 0; i < v.NumField(); i++ {
		columns = append(columns, strings.ToUpper(v.Type().Field(i).Name))
	}
	return columns
}

func (t *InstanceTable) Rows() []string {
	rows := []string{}
	for _, v := range t.Instances {
		rows = append(rows, v.String())
	}
	return rows
}

func (t *InstanceTable) OnTableResize(nNew int) {

	nOld := t.RowsDisplayed
	t.RowsDisplayed = nNew

	if nNew > nOld {
		if t.Offset > 0 {
			if t.Offset+t.Cursor > nNew {
				t.Offset = t.Offset + t.Cursor - nNew + 1
				t.Cursor = nNew - 1
			} else {
				t.Cursor += t.Offset
				t.Offset = 0
			}
		}
		return
	}

	if nNew < nOld {
		if t.Offset == 0 {
			if t.Cursor > nNew-1 {
				t.Offset += (t.Cursor + 1 - nNew)
				t.Cursor = nNew - 1
			}
			return
		}
		if t.Offset > 0 {
			if t.Cursor > 0 {
				cursorOld := t.Cursor
				t.Cursor -= (nOld - nNew)
				if t.Cursor < 0 {
					t.Cursor = 0
				}
				t.Offset += cursorOld - t.Cursor
			}
		}
	}

}

func (t *InstanceTable) SetTitle(title string) {
	t.Title = title
}

func (t *InstanceTable) SetCursor(n int) {
	t.Cursor = n
}

func (t *InstanceTable) SetOffset(n int) {
	t.Offset = n
}

func (t *InstanceTable) DefaultTitle(n int) string {
	return formatDefaultTitle(n)
}

func (t *InstanceTable) GetPositionByInstanceName(s string) (int, int) {
	rowsDisplayed := t.RowsDisplayed
	resIndex := 0
	s = strings.ToLower(s)
	found := false

	// TODO Evaluate the performance of the search by cycle on the list O(n)
	// and if is it not performant evaluate a tree or trie to get a logarithmic time complexity
	for i, v := range t.Instances {
		if strings.HasPrefix(strings.ToLower(v.Name), s) {
			resIndex = i
			found = true
			break
		}
	}
	if found {
		if resIndex < rowsDisplayed {
			return resIndex, 0
		}
		return (rowsDisplayed - 1), (resIndex - rowsDisplayed + 1)
	}
	return -1, 0
}

func (t *InstanceTable) Render(screen tcell.Screen, startY int) {
	columns := t.Columns()
	sw, sh := screen.Size()
	n := len(columns)
	delta := sw / n
	for delta < 21 {
		n--
		if n == 0 {
			return
		}
		delta = sw / n
	}
	w := 2
	DrawTableBox(screen, 0, startY-1, sw-1, sh-2)
	DrawStr(screen, sw/2-len(t.Title)/2-1, startY-1, tcell.StyleDefault, t.Title)

	for _, v := range columns[0:n] {
		DrawStr(screen, w, startY+1, styles[TopRow], v)
		w += delta
	}
	tableInstances := t.Instances[t.Offset:len(t.Instances)]

	for i, v := range tableInstances {
		w = 2
		if i+2+startY > sh-4 {
			return
		}
		targetStyle := styles[v.State]
		if i == t.Cursor {
			targetStyle = styles[SelectedRow]
		}
		for c, str := range strings.Split(v.String(), " ") {
			if c != n {
				DrawStr(screen, w, startY+i+2, targetStyle, str)
				// Fill gaps drawing blank chars
				for j := (w + len(str)); j < (w + delta); j++ {
					DrawStr(screen, j, startY+i+2, targetStyle, " ")
				}
				w += delta
			} else {
				break
			}
		}
	}
}

func NewInstanceTable(instances []client.Instance, n int) *InstanceTable {
	return &InstanceTable{
		Instances:     instances,
		Title:         formatDefaultTitle(n),
		Cursor:        0,
		Offset:        0,
		RowsDisplayed: n,
	}
}

func formatDefaultTitle(n int) string {
	return emoji.Sprintf(" :computer: EC2 Instances (%d) ", n)
}
