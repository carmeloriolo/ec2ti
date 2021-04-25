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

// func (t *InstanceTable) OnTableResize(numberOfRowsDisplaed int) {
//   oldRowsDisplayed := t.RowsDisplayed
//   t.RowsDisplayed = numberOfRowsDisplaed
//   if t.RowsDisplayed > 0 { // If there are rows to visualize
//
//     if t.RowsDisplayed < oldRowsDisplayed { // If the rows are less than before
//       if t.Cursor == 0 { // If the first element is selected
//         if t.Offset > 0 {
//           t.Offset--
//         }
//       } else if t.Cursor == oldRowsDisplayed-1 { // If the last element is selected
//         t.Offset += (oldRowsDisplayed - t.RowsDisplayed) // Offset is increased by difference
//         t.Cursor = t.RowsDisplayed - 1                   // Cursor always points to last element
//       } else { // If any other element in the middle is selected
//         if t.Offset > 0 { // Move sliding window
//           t.Offset += (oldRowsDisplayed - t.RowsDisplayed)
//           t.Cursor -= (oldRowsDisplayed - t.RowsDisplayed)
//         }
//         if t.Cursor > t.RowsDisplayed-1 {
//           t.Offset = t.Cursor - t.RowsDisplayed + 1
//           t.Cursor = t.RowsDisplayed - 1
//         }
//       }
//       return
//     }
//
//     if t.RowsDisplayed > oldRowsDisplayed { // If the rows are more than before
//       if t.Offset > 0 {
//         t.Cursor += t.Offset
//         t.Offset -= (t.RowsDisplayed - oldRowsDisplayed)
//       }
//     }
//
//     // Setting Lower bounds
//     if t.Cursor < 0 {
//       t.Cursor = 0
//     }
//     if t.Offset < 0 {
//       t.Offset = 0
//     }
//     // Setting Upper bounds
//     if t.Cursor > t.RowsDisplayed {
//       t.Cursor = t.RowsDisplayed - 1
//     }
//     if t.Offset > len(t.Instances)-1 {
//       t.Offset = len(t.Instances) - 1
//     }
//
//   } else {
//     // Reset table
//     t.Cursor = 0
//     t.Offset = 0
//   }
// }

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
