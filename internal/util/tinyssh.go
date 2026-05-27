package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs/appconfig"
	"time"
)

type uSSH struct{}

var SSH uSSH

func (u *uSSH) SSHRunRoot(config appconfig.SSHConfig, cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger.Info("cmd:", cmd)
	return u.SSHRunCommandOutputContext(ctx, config, u.sudoCommand(config, cmd))
}

func (u *uSSH) SSHRun(config appconfig.SSHConfig, cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger.Info("cmd:", cmd)
	return u.SSHRunCommandOutputContext(ctx, config, cmd)
}

func (u *uSSH) SCPUpload(config appconfig.SSHConfig, localFile string, remoteFile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return u.SCPUploadContext(ctx, config, localFile, remoteFile)
}

// ------------------------------------------------

func (u *uSSH) SCPUploadContext(ctx context.Context, config appconfig.SSHConfig, localFile string, remoteFile string) error {
	args := []string{}
	if config.Port > 0 {
		args = append(args, "-P", strconv.Itoa(config.Port))
	}
	args = append(args, "-o", "StrictHostKeyChecking=no")
	args = append(args, localFile, u.sshHost(config)+":"+remoteFile)
	_, err := u.runCommandContext(ctx, "scp", args...)
	return err
}

//func (u *uSSH) SSHRunCommandContext(ctx context.Context, config appconfig.SSHConfig, cmd string) error {
//	_, err := u.SSHRunCommandOutputContext(ctx, config, cmd)
//	return err
//}

func (u *uSSH) SSHRunCommandOutputContext(ctx context.Context, config appconfig.SSHConfig, cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)
	logger.Info("host: ", u.sshHost(config))
	//logger.Info("cmd: ", cmd)
	args := u.sshArgs(config, cmd)
	result, err := u.runCommandContext(ctx, "ssh", args...)
	logRes := result
	if len(logRes) > 100 {
		logRes = logRes[:100] + "..."
	}
	logger.Info("Output of SSH command:", logRes)
	if err != nil {
		logger.Error("Error executing SSH command:", err)
	}
	return result, err
}

func (u *uSSH) sudoCommand(config appconfig.SSHConfig, cmd string) string {
	return fmt.Sprintf("echo %s | sudo -S -p '' %s", config.Password, strings.TrimSpace(cmd))
}

func (u *uSSH) sshHost(config appconfig.SSHConfig) string {
	host := config.Host
	if config.User != "" {
		host = config.User + "@" + config.Host
	}
	return host
}

func (u *uSSH) sshArgs(config appconfig.SSHConfig, cmd string) []string {
	if config.Port == 0 {
		config.Port = 22
	}
	args := []string{}
	if config.Port > 0 {
		args = append(args, "-p", strconv.Itoa(config.Port))
	}
	args = append(args, "-o", "StrictHostKeyChecking=no")
	// args = append(args, "-o", "UserKnownHostsFile=/dev/null")
	args = append(args, u.sshHost(config), cmd)
	return args
}

func (u *uSSH) runCommandContext(ctx context.Context, name string, args ...string) (string, error) {
	var output bytes.Buffer
	command := exec.CommandContext(ctx, name, args...)
	//command.Stdout = io.MultiWriter(os.Stdout, &output)
	command.Stdout = &output
	command.Stderr = io.MultiWriter(os.Stderr, &output)
	err := command.Run()
	return output.String(), err
}
