package main

import (
	"fmt"
	"sucicada/home/internal/monitor"
	"sucicada/home/internal/util"
	"sucicada/home/internal/util/tinyssh"

	"github.com/gin-gonic/gin"
)

func main() {
	util.Server.InitExitSignal(func() {
		fmt.Println("exit")
		tinyssh.SSH.CloseMasters()
	})

	r := gin.Default()
	monitor.InitRouter(r)
	r.Run(":41407")
}
