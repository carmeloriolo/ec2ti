package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/carmeloriolo/ec2ti/internal/client"
	"github.com/carmeloriolo/ec2ti/internal/ui"
	"github.com/gdamore/tcell/v2"
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
		},
	}
)

func main() {
	app := &cli.App{
		Name:  appName,
		Usage: appDescription,
		Flags: appFlags,
		Action: func(c *cli.Context) error {
			cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(c.String(awsRegion)))
			if err != nil {
				log.Fatal(err)
			}
			instances, err := client.NewEc2Client(cfg).GetInstances()
			if err != nil {
				log.Fatal(err)
			}

			title := fmt.Sprintf(" EC2 Instances (%d) ", len(instances))
			u := ui.NewUi().SetTitle(title).SetTable(&ui.InstanceTable{
				Instances: instances,
			}).SetHandlers(ui.DefaultHandlers)

			for {
				switch ev := u.Screen.PollEvent().(type) {
				case *tcell.EventResize:
					u.ResetPosition()
					ui.Render(u)
				case *tcell.EventKey:
					if f, present := u.Handlers[ev.Key()]; present {
						f(u)
					}
				}
			}

		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
