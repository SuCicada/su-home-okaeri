package tinyssh

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs/appconfig"
	"sync"
	"syscall"
	"time"
)

// ------------------------------------------------
// 常驻 master 连接
//
// 每条命令单独 ssh 登录，在目标机上每次都会走一遍 PAM + logind 建会话，
// journal 里一次就是七八行。用 ControlMaster 复用同一条连接后，
// 只有 master 那一次会留下登录记录，后续开 channel 完全不写日志。
// 这里不用 ssh 自带的 ControlPersist 后台 master：那个 master 会被
// reparent 到 PID 1，退出后没人 wait 就变僵尸。改成由本进程直接持有。

const (
	// sshd 的 MaxSessions 默认是 10，留点余量，超了会 "Session open refused by peer"。
	sshMaxSessions = 8
	// master 建连（含认证）的等待上限。
	sshMasterReadyTimeout = 15 * time.Second
)

type sshMaster struct {
	lock   sync.Mutex
	cmd    *exec.Cmd
	cancel context.CancelFunc
	done   chan struct{}
	sem    chan struct{}
}

var (
	sshMastersLock sync.Mutex
	sshMasters     = map[string]*sshMaster{}

	sshControlPathOnce sync.Once
	sshControlPathVal  string
)

// sshControlPath 返回本程序专用的 control socket 路径模板。
// 不复用 ~/.ssh/config 里的路径，避免和用户交互式终端的 master 互相干扰。
// %C 由 ssh 自己展开成 (本机, 目标host, port, user) 的 hash，
// 顺便保证路径够短 —— unix socket 路径上限只有 108 字节。
func sshControlPath() string {
	sshControlPathOnce.Do(func() {
		base := os.Getenv("XDG_RUNTIME_DIR")
		if base == "" {
			base = os.TempDir()
		}
		dir := filepath.Join(base, "su-home-ssh")
		if err := os.MkdirAll(dir, 0o700); err != nil {
			logger.Error("创建 SSH control 目录失败: ", err)
			dir = os.TempDir()
		}
		sshControlPathVal = filepath.Join(dir, "cm-%C")
	})
	return sshControlPathVal
}

func (u *uSSH) master(config appconfig.SSHConfig) *sshMaster {
	key := fmt.Sprintf("%s:%d", u.sshHost(config), config.Port)

	sshMastersLock.Lock()
	defer sshMastersLock.Unlock()

	m, ok := sshMasters[key]
	if !ok {
		m = &sshMaster{sem: make(chan struct{}, sshMaxSessions)}
		sshMasters[key] = m
	}
	return m
}

// acquire 限制同时打开的 channel 数，别撞上 sshd 的 MaxSessions。
func (m *sshMaster) acquire(ctx context.Context) (func(), error) {
	select {
	case m.sem <- struct{}{}:
		return func() { <-m.sem }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (m *sshMaster) ensure(u *uSSH, config appconfig.SSHConfig) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.aliveLocked() && u.controlCommand(config, "check") == nil {
		return nil
	}
	m.stopLocked()
	// 可能有上次进程被强杀留下的死 socket，先清掉。
	u.controlCommand(config, "exit")
	return m.start(u, config)
}

// reset 让下一次调用重建 master。
func (m *sshMaster) reset(u *uSSH, config appconfig.SSHConfig) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.stopLocked()
	u.controlCommand(config, "exit")
}

func (m *sshMaster) aliveLocked() bool {
	if m.cmd == nil || m.done == nil {
		return false
	}
	select {
	case <-m.done:
		return false
	default:
		return true
	}
}

func (m *sshMaster) start(u *uSSH, config appconfig.SSHConfig) error {
	args := u.portArgs(config, "-p")
	args = append(args, u.commonArgs()...)
	args = append(args,
		"-M", "-N",
		"-o", "ControlPath="+sshControlPath(),
		// 不让 ssh 自己把 master 转后台，进程由本程序持有并回收。
		"-o", "ControlPersist=no",
		// 对端静默掉线时主动发现，而不是留一个连不上的死 socket。
		"-o", "ServerAliveInterval=30",
		"-o", "ServerAliveCountMax=3",
	)
	args = append(args, u.sshHost(config))

	ctx, cancel := context.WithCancel(context.Background())
	var stderr syncBuffer
	cmd := exec.CommandContext(ctx, "ssh", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	// 默认是 SIGKILL，master 来不及删 socket；给它机会正常收尾。
	cmd.Cancel = func() error { return cmd.Process.Signal(syscall.SIGTERM) }
	cmd.WaitDelay = 5 * time.Second

	if err := cmd.Start(); err != nil {
		cancel()
		return err
	}

	done := make(chan struct{})
	m.cmd, m.cancel, m.done = cmd, cancel, done
	go func() {
		err := cmd.Wait()
		close(done)
		if err != nil && ctx.Err() == nil {
			logger.Error("SSH master 退出: ", err)
		}
	}()

	deadline := time.Now().Add(sshMasterReadyTimeout)
	for time.Now().Before(deadline) {
		if u.controlCommand(config, "check") == nil {
			logger.Info("SSH master 已建立: ", u.sshHost(config))
			return nil
		}
		select {
		case <-done:
			m.clearLocked()
			return fmt.Errorf("ssh master 启动失败: %s", strings.TrimSpace(stderr.String()))
		case <-time.After(200 * time.Millisecond):
		}
	}

	m.stopLocked()
	return fmt.Errorf("ssh master 就绪超时: %s", u.sshHost(config))
}

func (m *sshMaster) stopLocked() {
	if m.cancel != nil {
		m.cancel()
	}
	if m.done != nil {
		select {
		case <-m.done:
		case <-time.After(6 * time.Second):
			logger.Error("SSH master 未在超时内退出")
		}
	}
	m.clearLocked()
}

func (m *sshMaster) clearLocked() {
	m.cmd, m.cancel, m.done = nil, nil, nil
}
