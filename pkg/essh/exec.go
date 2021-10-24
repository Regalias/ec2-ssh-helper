package essh

import (
	"fmt"
	"os"
	"os/exec"
)

type SshTarget struct {
	TargetHost   *SshHost
	IdentityFile string
	Bastions     []SshHost
}

type SshHost struct {
	Ip       string
	Port     uint16
	Username string
}

func (host *SshHost) getHostPortString(isHost bool) (opts string) {
	if host.Port != DEFAULT_SSH_PORT && !isHost {
		return fmt.Sprintf("%s@%s:%d", host.Username, host.Ip, host.Port)
	} else {
		return fmt.Sprintf("%s@%s", host.Username, host.Ip)
	}
}

func buildSshCommandOpts(target *SshTarget) []string {
	// Add identity file
	// -i <IDENTITY_FILE>
	sshOpts := []string{"-i", target.IdentityFile}

	if len(target.Bastions) > 0 {
		// Add bastion jump hosts
		// -J <USER@HOST:PORT> <USER@HOST:PORT>
		sshOpts = append(sshOpts, "-J")
		for _, bastion := range target.Bastions {
			sshOpts = append(sshOpts, bastion.getHostPortString(false))
		}
	}

	// Add target host
	// <USER@HOST>
	sshOpts = append(sshOpts, target.TargetHost.getHostPortString(true))

	// Add target port
	// -p <PORT>
	if target.TargetHost.Port != DEFAULT_SSH_PORT {
		sshOpts = append(sshOpts, "-p", fmt.Sprintf("%d", target.TargetHost.Port))
	}

	return sshOpts
}

func OpenSshInteractiveShell(target *SshTarget) error {

	cmd := exec.Command("ssh", buildSshCommandOpts(target)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("could not open SSH shell: %v", err)
	}

	return nil
}
