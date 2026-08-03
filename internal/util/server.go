package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type uServer struct {
}

var Server = &uServer{}

func (s *uServer) InitExitSignal(Close func()) {
	// システム信号を受信する
	go func() {
		sigChan := make(chan os.Signal, 1)
		defer close(sigChan)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		fmt.Printf("get signal: %v\n", sig)

		Close()

		os.Exit(0)
	}()
}
