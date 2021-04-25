package ui

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/carmeloriolo/ec2ti/internal/client"
)

type HeaderInterface interface {
	Rows() []string
}

type TableInterface interface {
	Columns() []string
	Rows() []string
	OnTableResize(int)
}

type InstanceTable struct {
	Instances     []client.Instance
	Cursor        int
	Offset        int
	RowsDisplayed int
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

	// Case1 . NRowsnew > NRowsold
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
		} else if t.Offset > 0 {
			t.Offset += (nOld - nNew)
			t.Cursor -= (nOld - nNew)
			if t.Cursor < 0 {
				t.Cursor = 0
			}
		}
	}

}

type InfoHeader struct {
	UserIdentity client.CallerIdentity
	Region       string
}

func (u *InfoHeader) Rows() []string {
	return []string{
		fmt.Sprintf("UserId:\t%s", u.UserIdentity.UserId),
		fmt.Sprintf("Account:\t%s", u.UserIdentity.Account),
		fmt.Sprintf("Arn:\t\t%s", u.UserIdentity.Arn),
		fmt.Sprintf("Region:\t%s", u.Region),
	}
}

func NewInstanceTable(instances []client.Instance, n int) *InstanceTable {

	return &InstanceTable{
		Instances:     instances,
		Cursor:        0,
		Offset:        0,
		RowsDisplayed: n,
	}

}
