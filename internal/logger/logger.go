package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugar       *zap.SugaredLogger
	projectRoot string
)

// func Sync() { _ = sugar.Sync() }
func Info(args ...any)  { sugar.Info(args...) }
func Debug(args ...any) { sugar.Debug(args...) }
func Warn(args ...any)  { sugar.Warn(args...) }
func Error(args ...any) { sugar.Error(args...) }

func init() {
	_, thisFile, _, _ := runtime.Caller(0)
	projectRoot = findProjectRoot(thisFile)

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	// encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	// encoderCfg.EncodeCaller = zapcore.FullCallerEncoder // ← 关键！改为完整路径
	encoderCfg.EncodeCaller = customCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Development(),
	)
	defer logger.Sync()
	sugar = logger.Sugar()
	//sugar.Info("logger init success")
}

// 从任意子文件路径向上找 go.mod 所在目录，即项目根目录
func findProjectRoot(file string) string {
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // 到达文件系统根，兜底
		}
		dir = parent
	}
	return filepath.Dir(file)
}

// 自定义 CallerEncoder，从项目根目录开始截取路径
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	fullPath := caller.FullPath()
	prefix := projectRoot + "/"
	if strings.HasPrefix(fullPath, prefix) {
		enc.AppendString(fullPath[len(prefix):])
		return
	}
	enc.AppendString(caller.TrimmedPath())
}
