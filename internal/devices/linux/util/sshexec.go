package linuxutil

import (
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/util"
)

func sshLinux(cmd string) (string, error) {
	return util.SSH.SSHRun(cfg.GetSSHConfig("linux"), cmd)
}
