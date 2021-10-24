package main

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pterm/pterm"
	"github.com/regalias/ec2-ssh-helper-go/pkg/essh"
	"github.com/regalias/ec2-ssh-helper-go/pkg/keygen"
)

func entrypoint(region string) {

	showtitle()

	pterm.Info.Printf("Using region '%s'\n", region)

	// Init config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, func(lo *config.LoadOptions) error {
		lo.Region = region
		return nil
	})
	checkFatalError(err)

	helper := essh.New(cfg)

	// Fetch EC2 instance list
	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(false).WithShowTimer(true).Start("[ec2:DescribeInstances] Fetching running EC2 instances...")
	check(err)
	instances, err := helper.GetRunningInstanceList(ctx)
	checkSpinnerError(spinner, err)

	if len(instances) < 1 {
		spinner.Warning(fmt.Sprintf("Found 0 running instances in '%s'. Aborting", region))
		return
	}
	spinner.Success(fmt.Sprintf("Found %d running instances in '%s'", len(instances), region))

	// Show EC2 instance details in table
	header := []string{"Instance Name", "Instance ID", "Type", "VPC", "Subnet", "AZ"}
	tableData := [][]string{header}

	var instanceOptList []string = make([]string, len(instances))

	for i, instance := range instances {
		instanceName := essh.GetInstanceName(&instance)
		tableData = append(tableData, []string{
			instanceName,
			*instance.InstanceId,
			string(instance.InstanceType),
			*instance.VpcId,
			*instance.SubnetId,
			*instance.Placement.AvailabilityZone,
		})
		instanceOptList[i] = fmt.Sprintf("%s (%s)", instanceName, *instance.InstanceId)
	}

	pterm.Println()
	check(pterm.DefaultTable.WithHasHeader(true).WithBoxed(false).WithData(tableData).Render())

	// Get target instance selection
	var targetInstanceIndex int
	err = survey.AskOne(&survey.Select{
		Message: "Select target instance: ",
		Options: instanceOptList,
	}, &targetInstanceIndex)

	if err != nil {
		if err == terminal.InterruptErr {
			pterm.Warning.Println("(Ctrl+C) Aborted by user, exiting")
			return
		}
		checkFatalError(err)
	}

	// For sanity, validate that it's in range
	if instanceOptList[targetInstanceIndex] == "" {
		pterm.Fatal.WithFatal().Print("Invalid input!")
	}
	pterm.Println()
	pterm.Info.Printf("Selected instance %s\n", *instances[targetInstanceIndex].InstanceId)

	genSpinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(false).WithShowTimer(true).Start("Generating temporary SSH session key...")
	check(err)
	// Create tempdir for temporary session SSH keypair
	temp, err := os.MkdirTemp("", "essh")
	checkSpinnerError(spinner, err)
	defer os.RemoveAll(temp)
	pterm.Debug.Printf("Creating temp dir: %s\n", temp)

	// Generate new key
	km := keygen.New(temp)
	privateKeyPath, pubKeyMaterial, err := km.GenerateKey()
	checkSpinnerError(spinner, err)
	pterm.Debug.Printf("Wrote key to: %s\n", privateKeyPath)

	genSpinner.Success("Generated temporary session SSH key")

	// Establish SSH session
	// TODO: this stuff
	username := "ec2-user"
	var port uint16 = 22

	// Install the temporary key
	keySpinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(false).WithShowTimer(true).Start("[ec2:SendSSHPublicKey] Installing temporary SSH session key...")
	check(err)
	err = helper.InstallKey(ctx, instances[targetInstanceIndex], username, pubKeyMaterial)
	checkSpinnerError(spinner, err)
	keySpinner.Success("Installed temporary SSH key")

	target := essh.GetTargetFromInstances(&instances[targetInstanceIndex], port, username, privateKeyPath, nil)

	pterm.Info.Printf("Opening SSH shell to %s@%s...\n", target.TargetHost.Username, target.TargetHost.Ip)

	pterm.DisableOutput()
	_, err = essh.OpenSshInteractiveShell(target)

	pterm.EnableOutput()
	checkFatalError(err)

	pterm.Info.Println("Session ended")
	// pterm.Info.Printf("To connect again, use '%s'\n", reconn)

}
