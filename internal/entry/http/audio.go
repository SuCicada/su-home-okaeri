package httpentry

import (
	"io"
	"mime/multipart"
	"os"
	"sucicada/home/internal/app/audio"
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
	defer cleanup()

	if err := audio.PlayAudio(ctx.Param("device"), command); err != nil {
		response.Error(ctx, err)
		return
	}
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
