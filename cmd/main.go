package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/carmeloriolo/ec2ti/internal/client"
	"github.com/carmeloriolo/ec2ti/internal/ui"
	"github.com/urfave/cli/v2"
)

const (
	appName        = "Ec2Ti"
	appDescription = "The terminal user interface to connect to your AWS EC2 instances easily"
	appVersion     = ""
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
		Name:    appName,
		Usage:   appDescription,
		Flags:   appFlags,
		Version: appVersion,
		Action: func(c *cli.Context) error {
			cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(c.String(awsRegion)))
			if err != nil {
				log.Fatal(err)
			}
			instances, err := client.NewEc2Client(cfg).GetInstances()
			if err != nil {
				log.Fatal(err)
			}
			userIdentity, err := client.NewSts(cfg).GetCallerIdentity()
			if err != nil {
				log.Fatal(err)
			}
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
