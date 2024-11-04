//go:build !windows

package localbinary

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func getCommand(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	gid := os.Getenv(PluginGID)
	uid := os.Getenv(PluginUID)
	if uid != "" && gid != "" {
		uid, err := strconv.Atoi(uid)
		if err != nil {
			return nil, fmt.Errorf("error parsing user ID: %w", err)
		}
		gid, err := strconv.Atoi(gid)
		if err != nil {
			return nil, fmt.Errorf("error parsing group ID: %w", err)
		}
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
			},
		}
	}
	return cmd, nil
}
