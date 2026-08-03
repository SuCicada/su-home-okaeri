package main

import (
	"fmt"
	"sucicada/home/internal/monitor"
	"sucicada/home/internal/util"

	"github.com/gin-gonic/gin"
)

func main() {
	util.Server.InitExitSignal(func() {
		fmt.Println("exit")
	})

	r := gin.Default()
	monitor.InitRouter(r)
	r.Run(":41407")
}
