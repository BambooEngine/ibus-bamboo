package introspectors

import (
	"fmt"
	"os/exec"
	"strings"
)

var introspectorCommand string = "<uninitialized command>"

func CmdSetIntrospectCommand(cmd string) {
	introspectorCommand = cmd
}

func CmdGetFocusedWmClass() string {
	cmd := exec.Command(introspectorCommand)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to get wmclass from command (%s)", introspectorCommand)
		return ""
	}
	return strings.Trim(string(out), "\r\n ")
}
