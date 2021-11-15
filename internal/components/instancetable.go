package components

import (
	"reflect"
	"strings"

	"github.com/carmeloriolo/ec2ti/internal/client"
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
