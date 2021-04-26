package client

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/goombaio/namegenerator"
)

func getRandomName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	return nameGenerator.Generate()
}

func getRandomType() string {
	min := 0
	max := 3
	n := rand.Intn(max-min) + min
	switch n {
	case 0:
		return "t2.micro"
	case 1:
		return "t3.large"
	case 2:
		return "m1.xlarge"
	default:
		return "t3.xlarge"
	}
}

func getRandomState() string {
	min := 0
	max := 5
	n := rand.Intn(max-min) + min
	return AllStates[n]
}

func GetMockedInstances(n int) []Instance {
	instances := []Instance{}
	for i := 0; i < n; i++ {
		name := getRandomName()
		id := time.Now().Unix() + int64(i)
		instances = append(instances, Instance{
			Id:           fmt.Sprintf("id-%d", id),
			Name:         name,
			Keyname:      fmt.Sprintf("key_%s", name),
			State:        getRandomState(),
			InstanceType: getRandomType(),
			Ip:           fmt.Sprintf("192.168.1.%d", i+10),
			LaunchTime:   "2021-04-26/15:04:05",
		})
	}
	return instances
}

func GetMockedUser() *CallerIdentity {
	return &CallerIdentity{
		UserId:  "ec2ti",
		Account: "ec2ti-account",
		Arn:     "arn:aws:iam::123456789012:user/ec2ti",
	}
}
