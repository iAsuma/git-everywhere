package qcmd

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gproc"
)

func ShellRun(pwd string, command string) error {
	return gproc.ShellRun(fmt.Sprintf("cd %s;%s", pwd, command))
}

func MustShellExec(pwd string, command string) string {
	return gproc.MustShellExec(fmt.Sprintf("cd %s;%s", pwd, command))
}
