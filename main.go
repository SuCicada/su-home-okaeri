package main

import (
	"SuCicada/home/internal"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := gin.Default()
	internal.GetRoute(r)
	r.Run(":41406")
}
