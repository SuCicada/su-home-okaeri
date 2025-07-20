package main

import (
	"SuCicada/home/internal"
	"SuCicada/home/internal/cfg"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	godotenv.Load()

	// 加载YAML配置文件
	cfg.LoadConfig("config.yaml")

	r := gin.Default()
	internal.GetRoute(r)
	r.Run(":41406")
}
