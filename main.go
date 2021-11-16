package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/carmeloriolo/ec2ti/internal/client"
	"github.com/carmeloriolo/ec2ti/internal/components"
	"github.com/carmeloriolo/ec2ti/internal/ui"
	"github.com/urfave/cli/v2"
)

const (
	appName         = "Ec2Ti"
	appDescription  = "The terminal user interface to connect to your AWS EC2 instances easily"
	flagAwsRegion   = "region"
	flagAwsEndpoint = "endpoint"
)

var (
	AppVersion = ""
	appFlags   = []cli.Flag{
		&cli.StringFlag{
			Name:    flagAwsEndpoint,
			EnvVars: []string{"AWS_EC2_CUSTOM_ENDPOINT"},
		},
		&cli.StringFlag{
			Name:    flagAwsRegion,
			EnvVars: []string{"AWS_DEFAULT_REGION"},
		},
	}
)

func main() {

	app := &cli.App{
		Name:    appName,
		Usage:   appDescription,
		Flags:   appFlags,
		Version: AppVersion,
		Action: func(c *cli.Context) error {
			cfg, err := config.LoadDefaultConfig(
				context.Background(),
				config.WithRegion(c.String(flagAwsRegion)),
				config.WithEndpointResolver(
					aws.EndpointResolverFunc(
						func(service, region string) (aws.Endpoint, error) {
							if c.String(flagAwsEndpoint) != "" {
								return aws.Endpoint{
									PartitionID:   "aws",
									URL:           c.String(flagAwsEndpoint),
									SigningRegion: c.String(flagAwsRegion),
								}, nil
							}
							return aws.Endpoint{}, &aws.EndpointNotFoundError{} // default fallback
						},
					),
				),
			)
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
			u := ui.NewUi().SetTitle(appName).SetHeader(&components.InfoHeader{
				UserIdentity: *userIdentity,
				Region:       c.String(flagAwsRegion),
			})
			u = u.SetTable(components.NewInstanceTable(instances, u.NumberOfRowsDisplayed()))
			return u.SetHandlers(ui.DefaultHandlers).Run()
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
