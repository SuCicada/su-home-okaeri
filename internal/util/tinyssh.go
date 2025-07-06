package util

import (
	"SuCicada/home/internal/logger"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type SSHConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func SSHRunRoot(config SSHConfig, cmd string) string {
	return runSSH(config, cmd, true)
}
func SSHRun(config SSHConfig, cmd string) string {
	return runSSH(config, cmd, false)
}

func runSSH(config SSHConfig, cmd string, sudo bool) string {
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
	args = append(args, host, cmd)
	output, err := exec.Command("ssh", args...).
		CombinedOutput()
	fmt.Println("Output of SSH command:", string(output))
	if err != nil {
		fmt.Println("Error executing SSH command:", err)
		return err.Error()
	}
	return string(output)
}
