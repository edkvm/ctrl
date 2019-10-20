// +build darwin

package invoke

import (
	"os/exec"
)

func buildCmd(execName string) *exec.Cmd {
	// TODO: Add context to cmd
	cmd := exec.Command(execName)

	return cmd
}