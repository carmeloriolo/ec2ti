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
}

type InstanceTable struct {
	Instances []client.Instance
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
