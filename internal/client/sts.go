package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Sts struct {
	client *sts.Client
}

func NewSts(cfg aws.Config) *Sts {
	s := Sts{
		client: sts.NewFromConfig(cfg),
	}
	return &s
}

func (s *Sts) GetCallerArn() (string, error) {
	identity, err := s.client.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return *identity.Arn, nil
}
