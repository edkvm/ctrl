// +build linux

package invoke

import (
	"os/exec"
)

func buildCmd(execName string) *exec.Cmd {
	// TODO: Add context to cmd
	cmd := exec.Command(execName)
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	//	Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	//}

	return cmd
}






