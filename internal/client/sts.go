package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type CallerIdentity struct {
	userId  string
	account string
	arn     string
}

type StsClientInterface interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

type StsClient struct {
	client StsClientInterface
}

func NewSts(cfg aws.Config) *StsClient {
	s := StsClient{
		client: sts.NewFromConfig(cfg),
	}
	return &s
}

func (s *StsClient) GetCallerIdentity() (*CallerIdentity, error) {
	identity, err := s.client.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	return &CallerIdentity{
		userId:  *identity.UserId,
		account: *identity.Account,
		arn:     *identity.Arn,
	}, nil
}
