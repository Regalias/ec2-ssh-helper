package essh

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect"
)

type Ec2SshHelper struct {
	ec2ConnectClient *ec2instanceconnect.Client
	ec2Client        *ec2.Client
}

func New(cfg aws.Config) *Ec2SshHelper {
	return &Ec2SshHelper{
		ec2ConnectClient: ec2instanceconnect.NewFromConfig(cfg),
		ec2Client:        ec2.NewFromConfig(cfg),
	}
}
