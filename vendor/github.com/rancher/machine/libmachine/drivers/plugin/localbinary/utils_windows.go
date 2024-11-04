//go:build windows

package localbinary

import "os/exec"

func getCommand(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	return cmd, nil
}
