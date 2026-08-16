package tinyssh

import (
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
	sudocmd := fmt.Sprintf("sudo %s", strings.TrimSpace(cmd))

	return u.SSHRunCommandOutputContext(ctx, config, sudocmd)
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

// CloseMasters 关闭所有常驻 master 连接，进程退出时调用。
func (u *uSSH) CloseMasters() {
	sshMastersLock.Lock()
	all := make([]*sshMaster, 0, len(sshMasters))
	for _, m := range sshMasters {
		all = append(all, m)
	}
	sshMastersLock.Unlock()

	for _, m := range all {
		m.lock.Lock()
		m.stopLocked()
		m.lock.Unlock()
	}
}

// ------------------------------------------------

func (u *uSSH) SCPUploadContext(ctx context.Context, config appconfig.SSHConfig, localFile string, remoteFile string) error {
	m := u.master(config)
	mux := m.ensure(u, config) == nil

	release, err := m.acquire(ctx)
	if err != nil {
		return err
	}
	defer release()

	stdout, stderr, err := u.runCommandContext(ctx, "scp", u.scpArgs(config, localFile, remoteFile, mux)...)
	if err != nil && mux && isSSHMuxError(stderr) {
		logger.Error("SCP 复用连接失效，重建后直连重试")
		m.reset(u, config)
		stdout, stderr, err = u.runCommandContext(ctx, "scp", u.scpArgs(config, localFile, remoteFile, false)...)
	}
	if err != nil {
		if detail := firstNonEmpty(stdout, stderr); detail != "" {
			return fmt.Errorf("%w: %s", err, detail)
		}
	}
	return err
}

func (u *uSSH) SSHRunCommandOutputContext(ctx context.Context, config appconfig.SSHConfig, cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)
	logger.Info("host: ", u.sshHost(config))

	m := u.master(config)
	mux := true
	if err := m.ensure(u, config); err != nil {
		// master 起不来不应该让功能不可用，退回到每条命令单独连接的老行为。
		logger.Error("SSH 复用连接不可用，降级为直连: ", err)
		mux = false
	}

	release, err := m.acquire(ctx)
	if err != nil {
		return "", err
	}
	defer release()

	stdout, stderr, err := u.runCommandContext(ctx, "ssh", u.sshArgs(config, cmd, mux)...)
	if err != nil && mux && isSSHMuxError(stderr) {
		logger.Error("SSH 复用连接失效，重建后直连重试")
		m.reset(u, config)
		stdout, stderr, err = u.runCommandContext(ctx, "ssh", u.sshArgs(config, cmd, false)...)
	}

	logRes := stdout
	if len(logRes) > 100 {
		logRes = logRes[:100] + "..."
	}
	logger.Info("Output of SSH command:", logRes)
	if err != nil {
		logger.Error("Error executing SSH command:", err)
		if detail := firstNonEmpty(stdout, stderr); detail != "" {
			err = fmt.Errorf("%w: %s", err, detail)
		}
	}
	return stdout, err
}

func (u *uSSH) sshHost(config appconfig.SSHConfig) string {
	host := config.Host
	if config.User != "" {
		host = config.User + "@" + config.Host
	}
	return host
}

// portArgs 只在配置里显式写了端口时才传 -p/-P，
// 否则交给 ~/.ssh/config 里该 Host 的 Port 决定。
func (u *uSSH) portArgs(config appconfig.SSHConfig, flag string) []string {
	if config.Port > 0 {
		return []string{flag, strconv.Itoa(config.Port)}
	}
	return nil
}

// muxArgs 复用 master 的连接；mux 为 false 时退回独立连接。
// ControlMaster=no 的语义是「不要成为 master，但可以复用已有 socket」。
func (u *uSSH) muxArgs(mux bool) []string {
	path := "none"
	if mux {
		path = sshControlPath()
	}
	return []string{"-o", "ControlMaster=no", "-o", "ControlPath=" + path}
}

func (u *uSSH) commonArgs() []string {
	return []string{
		"-o", "StrictHostKeyChecking=no",
		// 连不上时快速失败，而不是等外层 context 把进程 kill 掉。
		"-o", "ConnectTimeout=10",
		"-o", "BatchMode=yes",
	}
}

func (u *uSSH) sshArgs(config appconfig.SSHConfig, cmd string, mux bool) []string {
	args := u.portArgs(config, "-p")
	args = append(args, u.commonArgs()...)
	args = append(args, u.muxArgs(mux)...)
	args = append(args, u.sshHost(config), cmd)
	return args
}

func (u *uSSH) scpArgs(config appconfig.SSHConfig, localFile string, remoteFile string, mux bool) []string {
	args := u.portArgs(config, "-P")
	args = append(args, u.commonArgs()...)
	args = append(args, u.muxArgs(mux)...)
	args = append(args, localFile, u.sshHost(config)+":"+remoteFile)
	return args
}

func (u *uSSH) runCommandContext(ctx context.Context, name string, args ...string) (string, string, error) {
	var output syncBuffer
	command := exec.CommandContext(ctx, name, args...)
	command.Stdout = io.MultiWriter(os.Stderr, &output)

	var errOutput syncBuffer
	command.Stderr = io.MultiWriter(os.Stderr, &errOutput)

	err := command.Run()
	if errOutput.Len() > 0 {
		logger.Error("Error output of SSH command:", errOutput.String())
	}
	return output.String(), errOutput.String(), err
}

// controlCommand 执行 ssh -O <action>，用于探活和清理。
// check 失败是正常路径（master 还没起来），所以这里不打日志。
func (u *uSSH) controlCommand(config appconfig.SSHConfig, action string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := u.portArgs(config, "-p")
	args = append(args, "-o", "ControlPath="+sshControlPath())
	args = append(args, "-O", action, u.sshHost(config))

	cmd := exec.CommandContext(ctx, "ssh", args...)
	return cmd.Run()
}

// isSSHMuxError 判断失败是不是复用连接本身的问题，
// 是的话值得重建 master 再用独立连接重试一次。
func isSSHMuxError(stderr string) bool {
	s := strings.ToLower(stderr)
	for _, sign := range []string{
		"session open refused by peer",
		"control socket",
		"controlsocket",
		"mux_client",
		"multiplex",
	} {
		if strings.Contains(s, sign) {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
