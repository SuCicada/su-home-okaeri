package httpentry

import (
	"sucicada/home/internal/app/media"
	"sucicada/home/internal/response"
	"sucicada/home/internal/structs"

	"github.com/gin-gonic/gin"
)

func RegisterMediaRoutes(r *gin.Engine) {
	r.GET("/api/devices/:device/media/status", getMediaStatus)
	r.POST("/api/devices/:device/media/command", executeMediaCommand)
}

func getMediaStatus(ctx *gin.Context) {
	noThumbnail := ctx.Query("nothumbnail")

	status, err := media.GetStatus(ctx.Param("device"))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if noThumbnail == "true" {
		status.Thumbnail = nil
	}
	response.Success(ctx, status)
}

func executeMediaCommand(ctx *gin.Context) {
	var command structs.MediaCommand
	if err := ctx.ShouldBindJSON(&command); err != nil {
		response.Bad(ctx, err.Error())
		return
	}
	if err := media.Execute(ctx.Param("device"), command); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
