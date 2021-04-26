package client

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var (
	AllStates = []string{
		"pending",
		"running",
		"stopping",
		"stopped",
		"shutting-down",
		"terminated",
	}
)

type Instance struct {
	Id           string
	Name         string
	State        string
	InstanceType string
	Keyname      string
	Ip           string
	LaunchTime   string
}

type ByName []Instance

func (b ByName) Len() int {
	return len(b)
}

func (b ByName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByName) Less(i, j int) bool {
	return b[i].Name < b[j].Name
}

func (i *Instance) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s", i.Id, i.Name, i.State, i.InstanceType, i.Keyname, i.Ip, i.LaunchTime)
}

type Ec2ClientInterface interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

type Ec2Client struct {
	client Ec2ClientInterface
}

func NewEc2Client(cfg aws.Config) *Ec2Client {
	return &Ec2Client{
		client: ec2.NewFromConfig(cfg),
	}
}

func (e *Ec2Client) GetInstances() ([]Instance, error) {
	return e.GetInstancesByState(AllStates)
}

func (e *Ec2Client) GetInstancesByState(states []string) ([]Instance, error) {

	instances := []Instance{}

	result, err := e.client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: states,
			},
		},
	})

	if err != nil {
		return nil, err
	}
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					keyname := ""
					if i.KeyName != nil {
						keyname = *i.KeyName
					}
					instances = append(instances, Instance{
						Name:         *t.Value,
						Id:           *i.InstanceId,
						State:        string(i.State.Name),
						InstanceType: string(i.InstanceType),
						Keyname:      keyname,
						Ip:           *i.PrivateIpAddress,
						LaunchTime:   i.LaunchTime.Format("2006-01-02/15:04:05"),
					})
				}
			}
		}

	}
	sort.Sort(ByName(instances))
	return instances, nil
}
