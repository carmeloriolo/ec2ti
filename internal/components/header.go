package components

import (
	"fmt"

	"github.com/carmeloriolo/ec2ti/internal/client"
)

type Header interface {
	Rows() []string
}
type InfoHeader struct {
	UserIdentity client.CallerIdentity
	Region       string
}

func (u *InfoHeader) Rows() []string {
	return []string{
		fmt.Sprintf("UserId: %s", u.UserIdentity.UserId),
		fmt.Sprintf("Account: %s", u.UserIdentity.Account),
		fmt.Sprintf("Arn: %s", u.UserIdentity.Arn),
		fmt.Sprintf("Region: %s", u.Region),
	}

}
