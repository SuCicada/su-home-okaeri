package main

import (
	"SuCicada/home/internal"
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/mqttpkg"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	mqttpkg.Init()
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

	// システム信号を受信する
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("get signal: %v\n", sig)

		Close()

		os.Exit(0)
	}()

	r := gin.Default()
	internal.GetRoute(r)
	r.Run(":41406")
}
