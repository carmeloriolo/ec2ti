package client

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

const (
	noInstances    = "noInstances"
	onlyRunning    = "onlyRunning"
	runningStopped = "runningStopped"
	ec2Error       = "ec2Error"
)

var (
	launchTime       = time.Date(2006, 01, 02, 15, 04, 05, 00, time.UTC).Format("2006-01-02/15:04:05")
	runningInstances = []Instance{
		{
			Name:         "ec2mock",
			Id:           "i-123456",
			State:        "running",
			InstanceType: "t2.micro",
			Keyname:      "keypair",
			Ip:           "192.168.1.1",
			LaunchTime:   launchTime,
		},
	}
	runningStoppedInstances = []Instance{
		{
			Name:         "ec2mock",
			Id:           "i-123456",
			State:        "running",
			InstanceType: "t2.micro",
			Keyname:      "keypair",
			Ip:           "192.168.1.1",
			LaunchTime:   launchTime,
		},
		{
			Name:         "ec2mock",
			Id:           "i-123456",
			State:        "stopped",
			InstanceType: "t2.micro",
			Keyname:      "keypair",
			Ip:           "192.168.1.2",
			LaunchTime:   launchTime,
		},
	}
)

type Ec2ClientMock struct {
	testCase string
}

func (e *Ec2ClientMock) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	switch e.testCase {
	case onlyRunning:
		ec2InstancesOutput := ec2.DescribeInstancesOutput{}
		stub, err := ioutil.ReadFile("./stubs/instances-running.json")
		if err != nil {
			log.Fatal(err)
		}
		_ = json.Unmarshal([]byte(string(stub)), &ec2InstancesOutput)
		return &ec2InstancesOutput, nil
	case runningStopped:
		ec2InstancesOutput := ec2.DescribeInstancesOutput{}
		stub, err := ioutil.ReadFile("./stubs/instances-running-stopped.json")
		if err != nil {
			log.Fatal(err)
		}
		_ = json.Unmarshal([]byte(string(stub)), &ec2InstancesOutput)
		return &ec2InstancesOutput, nil
	case noInstances:
		return &ec2.DescribeInstancesOutput{}, nil
	default:
		return nil, errors.New(ec2Error)
	}
}

func getEc2ClientMock(s string) *Ec2Client {
	return &Ec2Client{
		client: &Ec2ClientMock{testCase: s},
	}
}

func Test_GetInstancesByState(t *testing.T) {

	tests := []struct {
		name          string
		ec2           *Ec2Client
		states        []string
		wantInstances []Instance
		wantErr       error
	}{
		{
			name:          "Test no instances",
			ec2:           getEc2ClientMock(noInstances),
			states:        []string{"running"},
			wantInstances: []Instance{},
			wantErr:       nil,
		},
		{
			name:          "Test error retrieving instances",
			ec2:           getEc2ClientMock(ec2Error),
			states:        []string{"running"},
			wantInstances: nil,
			wantErr:       errors.New(ec2Error),
		},
		{
			name:          "Test running instances",
			ec2:           getEc2ClientMock(onlyRunning),
			states:        []string{"running", "stopped"},
			wantInstances: runningInstances,
			wantErr:       nil,
		},
		{
			name:          "Test running stopped instances",
			ec2:           getEc2ClientMock(runningStopped),
			states:        []string{"running", "stopped"},
			wantInstances: runningStoppedInstances,
			wantErr:       nil,
		},
	}
	for _, tt := range tests {
		st := tt
		t.Run(tt.name, func(t *testing.T) {
			instances, err := st.ec2.GetInstancesByState(allStates)
			if (st.wantErr == nil && err != nil) || (st.wantErr != nil && err == nil) {
				t.Errorf("GetInstancesByState(%s) st.wantErr = %s err = %s", st.states, st.wantErr, err)
				return
			}
			if st.wantInstances != nil && instances == nil {
				t.Errorf("GetInstancesByState(%s) st.wantInstances = %s instances = %s", st.states, st.wantInstances, instances)
				return
			}
			if !reflect.DeepEqual(st.wantInstances, instances) {
				t.Errorf("GetInstancesByState(%s) st.wantInstances = %s instances = %s", st.states, st.wantInstances, instances)
				return
			}
		})
	}

}
