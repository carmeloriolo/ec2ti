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

func (t *InstanceTable) OnTableResize(numberOfRowsDisplaed int) {
	oldRowsDisplayed := t.RowsDisplayed
	t.RowsDisplayed = numberOfRowsDisplaed
	if t.RowsDisplayed > 0 { // If there are rows to visualize

		if t.RowsDisplayed < oldRowsDisplayed { // If the rows are less than before
			if t.Cursor == 0 { // If the first element is selected
				if t.Offset > 0 {
					t.Offset--
				}
			} else if t.Cursor == oldRowsDisplayed-1 { // If the last element is selected
				t.Offset += (oldRowsDisplayed - t.RowsDisplayed) // Offset is increased by difference
				t.Cursor = t.RowsDisplayed - 1                   // Cursor always points to last element
			} else { // If any other element in the middle is selected
				if t.Offset > 0 { // Move sliding window
					t.Offset += (oldRowsDisplayed - t.RowsDisplayed)
					t.Cursor -= (oldRowsDisplayed - t.RowsDisplayed)
				}
				if t.Cursor > t.RowsDisplayed-1 {
					t.Offset = t.Cursor - t.RowsDisplayed + 1
					t.Cursor = t.RowsDisplayed - 1
				}
			}
			return
		}

		if t.RowsDisplayed > oldRowsDisplayed { // If the rows are more than before
			if t.Offset > 0 {
				t.Cursor += t.Offset
				if t.Cursor < 0 {
					t.Cursor = 0
				}
				t.Offset -= (t.RowsDisplayed - oldRowsDisplayed)
				if t.Offset < 0 {
					t.Offset = 0
				}
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
