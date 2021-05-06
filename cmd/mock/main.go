package main

import (
	"log"
	"os"

	"github.com/carmeloriolo/ec2ti/internal/client"
	"github.com/carmeloriolo/ec2ti/internal/ui"
	"github.com/urfave/cli/v2"
)

const (
	appName        = "Ec2Ti"
	appDescription = "The terminal user interface to connect to your AWS EC2 instances easily"
	awsRegion      = "region"
)

var (
	appFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    awsRegion,
			EnvVars: []string{"AWS_DEFAULT_REGION"},
			Value:   "eu-west-1",
		},
	}
)

func main() {
	app := &cli.App{
		Name:  appName,
		Usage: appDescription,
		Flags: appFlags,
		Action: func(c *cli.Context) error {
			instances := client.GetMockedInstances(64)
			userIdentity := client.GetMockedUser()
			u := ui.NewUi().SetTitle(appName).SetHeader(&ui.InfoHeader{
				UserIdentity: *userIdentity,
				Region:       c.String(awsRegion),
			})
			u = u.SetTable(ui.NewInstanceTable(instances, u.NumberOfRowsDisplayed()))
			return u.SetHandlers(ui.DefaultHandlers).Run()
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
