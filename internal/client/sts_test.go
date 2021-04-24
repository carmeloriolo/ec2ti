package client

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const (
	callerIdentity = "callerIdentity"
	errorIdentity  = "errorIdentity"
)

type StsClientMock struct {
	testCase string
}

func (s *StsClientMock) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	switch s.testCase {
	case callerIdentity:
		stsCallerIdentityOutput := sts.GetCallerIdentityOutput{}
		stub, err := ioutil.ReadFile("./stubs/get-caller-identity.json")
		if err != nil {
			log.Fatal(err)
		}
		_ = json.Unmarshal([]byte(string(stub)), &stsCallerIdentityOutput)
		return &stsCallerIdentityOutput, nil
	default:
		return nil, errors.New("error")
	}
}

func getStsMock(s string) *StsClient {
	return &StsClient{
		client: &StsClientMock{testCase: s},
	}
}

func Test_GetCallerIdentity(t *testing.T) {

	tests := []struct {
		name         string
		sts          *StsClient
		wantIdentity *CallerIdentity
		wantError    error
	}{
		{
			name:         "Testing error case",
			sts:          getStsMock(errorIdentity),
			wantIdentity: nil,
			wantError:    errors.New("error"),
		},
		{
			name: "Testing get caller identity",
			sts:  getStsMock(callerIdentity),
			wantIdentity: &CallerIdentity{
				UserId:  "AKIDSJHUFAHUDASMK",
				Account: "48120983120391",
				Arn:     "arn:aws:iam::48120983120391:user/user.mock",
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		st := tt
		t.Run(tt.name, func(t *testing.T) {
			identity, err := st.sts.GetCallerIdentity()
			if (st.wantError != nil && err == nil) || (st.wantError == nil && err != nil) {
				t.Errorf("GetCallerIdentity() st.wantError = %s err = %s", st.wantError, err)
				return
			}
			if !reflect.DeepEqual(st.wantIdentity, identity) {
				t.Errorf("GetCallerIdentity() st.wantIdentity = %s identity = %s", st.wantError, err)
				return
			}
		})
	}

}
