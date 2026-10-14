package main

import (
	"sucicada/home/internal"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/devices"
	mqttentry "sucicada/home/internal/entry/mqtt"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/mqttpkg"
	"sucicada/home/internal/util"
	"sucicada/home/internal/util/tinyssh"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	devices.Init()
	mqttpkg.Init()

	mqttRouter := mqttpkg.NewRouter()
	mqttentry.RegisterRoutes(mqttRouter)
	mqttpkg.UseRoutes(mqttRouter)
}
func InitHttp() {
	r := gin.Default()
	internal.GetRoute(r)
	r.Run(":41406")
}
func Close() {
	mqttpkg.Close()
	tinyssh.SSH.CloseMasters()
}

func main() {
	// 加载环境变量
	godotenv.Load()
	// 加载YAML配置文件
	logger.Info("Loading config file: ", cfg.GetConfigPath())
	cfg.LoadConfig(cfg.GetConfigPath())

	Init()
	util.Server.InitExitSignal(Close)
	InitHttp()
}
