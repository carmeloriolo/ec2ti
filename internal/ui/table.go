package ui

import (
	"reflect"
	"strings"

	"github.com/carmeloriolo/ec2ti/internal/client"
)

type Table interface {
	Columns() []string
	Rows() []string
}

type InstanceTable struct {
	Instances []client.Instance
}

func (t *InstanceTable) Columns() []string {
	headers := []string{}
	v := reflect.ValueOf(client.Instance{})
	for i := 0; i < v.NumField(); i++ {
		headers = append(headers, strings.ToUpper(v.Type().Field(i).Name))
	}
	return headers
}

func (t *InstanceTable) Rows() []string {
	rows := []string{}
	for _, v := range t.Instances {
		rows = append(rows, v.String())
	}
	return rows
}
