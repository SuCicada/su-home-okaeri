package httpentry

import (
	"io"
	"mime/multipart"
	"os"
	"sucicada/home/internal/app/audio"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/response"
	"sucicada/home/internal/structs"

	"github.com/gin-gonic/gin"
)

func RegisterAudioRoutes(r *gin.Engine) {
	r.POST("/api/devices/:device/audio/play", playAudio)
}

func playAudio(ctx *gin.Context) {
	command, cleanup, err := bindAudioPlayRequest(ctx)
	if err != nil {
		response.Bad(ctx, err.Error())
		return
	}

	// gin.Context 在 handler 返回后会被复用，goroutine 里不能再碰它。
	device := ctx.Param("device")

	// 播放会一直阻塞到音频放完，不能占着 HTTP 连接不放，
	// 否则调用方会先超时断开，然后重试，结果是两个 ffplay 叠着放。
	// cleanup 的所有权一并交给这个 goroutine —— handler 返回时不能删掉还没上传的文件。
	go func() {
		defer cleanup()
		if err := audio.PlayAudio(device, command); err != nil {
			logger.Error("play audio error: ", err)
		}
	}()

	response.Success(ctx)
}

func bindAudioPlayRequest(ctx *gin.Context) (structs.AudioPlayRequest, func(), error) {
	file, err := ctx.FormFile("audiofile")
	if err == nil {
		localFile, cleanup, err := saveUploadedAudioFile(file)
		if err != nil {
			return structs.AudioPlayRequest{}, func() {}, err
		}
		return structs.AudioPlayRequest{AudioFile: localFile}, cleanup, nil
	}

	if audioBase64 := ctx.PostForm("audiobase64"); audioBase64 != "" {
		return structs.AudioPlayRequest{AudioBase64: audioBase64}, func() {}, nil
	}

	var command structs.AudioPlayRequest
	if err := ctx.ShouldBindJSON(&command); err != nil {
		return structs.AudioPlayRequest{}, func() {}, err
	}
	return command, func() {}, nil
}

func saveUploadedAudioFile(fileHeader *multipart.FileHeader) (string, func(), error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", func() {}, err
	}
	defer src.Close()

	dst, err := os.CreateTemp("", "su-home-upload-audio-*.wav")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() {
		dst.Close()
		os.Remove(dst.Name())
	}

	if _, err := io.Copy(dst, src); err != nil {
		cleanup()
		return "", func() {}, err
	}
	if err := dst.Close(); err != nil {
		cleanup()
		return "", func() {}, err
	}

	return dst.Name(), cleanup, nil
}
