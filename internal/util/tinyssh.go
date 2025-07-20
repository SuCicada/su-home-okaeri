package util

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/logger"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func SSHRunRoot(config cfg.SSHConfig, cmd string) (string, error) {
	return runSSH(config, cmd, true)
}

func SSHRun(config cfg.SSHConfig, cmd string) (string, error) {
	return runSSH(config, cmd, false)
}

func runSSH(config cfg.SSHConfig, cmd string, sudo bool) (string, error) {
	if config.Port == 0 {
		config.Port = 22
	}
	host := config.Host
	if config.User != "" {
		host = config.User + "@" + config.Host
	}
	cmd = strings.TrimSpace(cmd)
	logger.Info("host:", host)
	logger.Info("cmd:", cmd)
	if sudo {
		cmd = fmt.Sprintf("echo %s | sudo -S -p '' %s ", config.Password, cmd)
	}

	args := []string{}
	if config.Port > 0 {
		args = append(args, "-p", strconv.Itoa(config.Port))
	}
	args = append(args, "-o", "StrictHostKeyChecking=no")
	// args = append(args, "-o", "UserKnownHostsFile=/dev/null")
	args = append(args, host, cmd)
	output, err := exec.Command("ssh", args...).
		CombinedOutput()
	var result = string(output)
	logger.Info("Output of SSH command:", result)
	if err != nil {
		logger.Error("Error executing SSH command:", err)
	}
	return result, err
}
