package main

import (
	"context"
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
			userIdentity, err := client.NewSts(cfg).GetCallerIdentity()
			if err != nil {
				log.Fatal(err)
			}
			// instances := []client.Instance{}
			// for i := 0; i < 64; i++ {
			//   instances = append(instances, client.Instance{
			//     Id:           fmt.Sprintf("id-%d", i),
			//     Name:         "ec2mock",
			//     State:        "running",
			//     InstanceType: "t2.micro",
			//     Ip:           "192.168.1.1",
			//     LaunchTime:   "jdaisodjaso",
			//   })
			// }
			// userIdentity := &client.CallerIdentity{
			//   UserId:  "carmelo",
			//   Account: "account",
			//   Arn:     "adshuèdhaudhasuarn",
			// }

			u := ui.NewUi()
			rows := u.NumberOfRowsDisplayed()
			u.SetTitle(appName).SetHeader(&ui.InfoHeader{
				UserIdentity: *userIdentity,
				Region:       c.String(awsRegion),
				// Region: "eu-west-1",
			}).SetTable(ui.NewInstanceTable(instances, rows)).SetHandlers(ui.DefaultHandlers)

			for {
				switch ev := u.Screen.PollEvent().(type) {
				case *tcell.EventResize:
					u.Render()
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
