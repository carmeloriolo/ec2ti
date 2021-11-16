package components

import (
	"fmt"

	"github.com/carmeloriolo/ec2ti/internal/client"
	"github.com/gdamore/tcell/v2"
	"github.com/kyokomi/emoji"
)

var (
	commandLabels = []string{
		searchLabel,
		describeLabel,
		shellLabel,
		ctrlCLabel,
	}
)

type Header interface {
	Render(string, tcell.Screen, int)
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

func (h *InfoHeader) Render(title string, screen tcell.Screen, endY int) {
	sw, _ := screen.Size()
	title = formatTitle(title)
	DrawHeaderBox(screen, 0, 0, sw-1, endY)
	DrawStr(screen, sw/2-len(title)/2-1, 0, styles[HeaderRow], title)
	nRows := endY - 3
	for i, r := range h.Rows() {
		if i < nRows {
			DrawStr(screen, 2, 2+i, styles[HeaderRow], r)
		} else {
			break
		}
	}
	if sw > 50 {
		for i, l := range commandLabels {
			if i < nRows {
				DrawStr(screen, sw-2-len(l), i+2, styles[CommandRow], l)
			} else {
				return
			}
		}
	}
}

func formatTitle(t string) string {
	return emoji.Sprintf(" :rocket: %s :beer:", t)
}
