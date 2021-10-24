package essh

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect"
)

func (h *Ec2SshHelper) InstallKey(ctx context.Context, instance types.Instance, username string, pubKeyData string) error {

	_, err := h.ec2ConnectClient.SendSSHPublicKey(ctx, &ec2instanceconnect.SendSSHPublicKeyInput{
		InstanceId:       instance.InstanceId,
		AvailabilityZone: instance.Placement.AvailabilityZone,
		InstanceOSUser:   &username,
		SSHPublicKey:     &pubKeyData,
	})
	if err != nil {
		return fmt.Errorf("failed to install SSH key: %v", err)
	}

	return nil
}
