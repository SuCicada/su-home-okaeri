package main

import (
	"fmt"
	"os"
	"os/signal"
	"sucicada/home/internal"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/devices"
	mqttentry "sucicada/home/internal/entry/mqtt"
	"sucicada/home/internal/mqttpkg"
	"sucicada/home/internal/service/mqttservice"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	devices.Init()
	mqttpkg.Init()
	mqttservice.RegisterRoutes()
	mqttentry.RegisterMediaRoutes()
}
func InitHttp() {
	r := gin.Default()
	internal.GetRoute(r)
	r.Run(":41406")
}
func Close() {
	mqttpkg.Close()
}

func main() {
	// 加载环境变量
	godotenv.Load()
	// 加载YAML配置文件
	cfg.LoadConfig("config.yaml")

	Init()
	InitExitSignal()
	InitHttp()
}

func InitExitSignal() {
	// システム信号を受信する
	go func() {
		sigChan := make(chan os.Signal, 1)
		defer close(sigChan)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		fmt.Printf("get signal: %v\n", sig)

		Close()

		os.Exit(0)
	}()
}
