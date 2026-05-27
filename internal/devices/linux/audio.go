package linux

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs"
	"sucicada/home/internal/util"
)

type sLinuxAudio struct{}

func (l *sLinuxAudio) Get() (any, error) {
	return nil, nil
}

func (l *sLinuxAudio) Set(command string) error {
	return nil
}

func (l *sLinuxAudio) PlayAudio(command structs.AudioPlayRequest) error {
	localFile, cleanup, err := prepareAudioFile(command)
	if err != nil {
		return err
	}
	defer cleanup()

	remoteFile := filepath.Join("/tmp", filepath.Base(localFile))
	sshConfig := cfg.GetSSHConfig(Device.Name)

	if err := util.SSH.SCPUpload(sshConfig, localFile, remoteFile); err != nil {
		logger.Error("scp upload error", err)
		return err
	}

	_, err = util.SSH.SSHRun(sshConfig, "ffplay -nodisp -autoexit "+remoteFile)
	return err
}

func prepareAudioFile(command structs.AudioPlayRequest) (string, func(), error) {
	if command.AudioFile != "" {
		return command.AudioFile, func() {}, nil
	}

	if command.AudioBase64 == "" {
		return "", func() {}, errors.New("audiofile or audiobase64 is required")
	}

	audioBytes, err := base64.StdEncoding.DecodeString(command.AudioBase64)
	if err != nil {
		return "", func() {}, err
	}

	localFile, err := os.CreateTemp("", "su-home-audio-*.wav")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() {
		os.Remove(localFile.Name())
		localFile.Close()
	}

	if _, err := localFile.Write(audioBytes); err != nil {
		logger.Error("scp upload error", err)
		cleanup()
		return "", func() {}, err
	}
	if err := localFile.Close(); err != nil {
		logger.Error("scp upload error", err)
		cleanup()
		return "", func() {}, err
	}

	return localFile.Name(), cleanup, nil
}
