package test

import (
	"fmt"
	"os"
	"testing"

	"SuCicada/home/service"

	"github.com/joho/godotenv"
	"github.com/melbahja/goph"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	os.Exit(m.Run())
}
func TestSleep(t *testing.T) {
	service.ChangeLightLinux(10)
}

func TestSsh(t *testing.T) {

	// auth, err := goph.UseAgent()
	// if err != nil {
	// 	// handle error
	// }

	client, err := goph.NewConn(&goph.Config{
		User: os.Getenv("LINUX_USER"),
		Addr: os.Getenv("LINUX_HOST"),
		Port: 22,
		// Auth: auth,
		// Timeout:  DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		assert.Nil(t, err)
	}

	output, err := client.Run("ls -la")
	if err != nil {
		assert.Nil(t, err)
	}
	fmt.Println(string(output))
}
