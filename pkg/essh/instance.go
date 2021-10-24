package essh

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (h *Ec2SshHelper) GetRunningInstanceList(ctx context.Context) (instanceList []types.Instance, err error) {

	paginator := ec2.NewDescribeInstancesPaginator(h.ec2Client, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %v", err)
		}
		for _, reservation := range page.Reservations {
			instanceList = append(instanceList, reservation.Instances...)
		}
	}

	return instanceList, nil
}

func GetTargetFromInstances(targetInstance *types.Instance, sshPort uint16, username string, privateKeyPath string, bastionInstance *types.Instance) *SshTarget {

	target := &SshTarget{
		TargetHost: &SshHost{
			Port:     sshPort,
			Username: username,
		},
		IdentityFile: privateKeyPath,
	}
	if bastionInstance != nil {
		// Add bastion host
		target.Bastions = []*SshHost{
			{
				Ip:       *bastionInstance.PublicIpAddress,
				Port:     sshPort,
				Username: username,
			},
		}
		// Use private IP for target host
		target.TargetHost.Ip = *targetInstance.PrivateIpAddress
	} else {
		// Use public IP for target host
		target.TargetHost.Ip = *targetInstance.PublicIpAddress
	}
	return target
}
