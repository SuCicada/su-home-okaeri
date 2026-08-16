package linuxutil

import (
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/util/tinyssh"
)

func sshLinux(cmd string) (string, error) {
	return tinyssh.SSH.SSHRun(cfg.GetSSHConfig("linux"), cmd)
}
