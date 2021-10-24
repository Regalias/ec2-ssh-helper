package essh

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (h *Ec2SshHelper) GetInstanceList(ctx context.Context) (instanceList []types.Instance, err error) {

	paginator := ec2.NewDescribeInstancesPaginator(h.ec2Client, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running "},
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
