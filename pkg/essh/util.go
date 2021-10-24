package essh

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetInstanceName(instance *types.Instance) (name string) {
	for _, tag := range instance.Tags {
		if strings.ToLower(*tag.Key) == "name" {
			return *tag.Value
		}
	}
	// No name, return '-'
	return "-"
}
