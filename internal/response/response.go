package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Success(ctx *gin.Context, data ...any) {
	if len(data) > 0 {
		ctx.JSON(http.StatusOK, Response{Message: "success", Data: data[0]})
	} else {
		ctx.JSON(http.StatusOK, Response{Message: "success"})
	}
}
func Bad(ctx *gin.Context, data any) {
	log.Println("response bad request: ", data)
	ctx.JSON(http.StatusBadRequest, Response{Message: "bad request", Data: data})
}
func Error(ctx *gin.Context, err error) {
	log.Println("response error: ", err)
	ctx.JSON(http.StatusInternalServerError,
		Response{Message: "error", Error: err.Error()})
}
