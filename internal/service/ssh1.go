package service

import (
	"bytes"
	"fmt"
	"log"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

func ssh2() {
	// auth, err := goph.UseAgent()
	// if err != nil {
	// 	// handle error
	// }

	client, err := goph.NewConn(&goph.Config{
		User: "root",
		Addr: "192.1.1.3",
		Port: 22,
		// Auth: auth,
		// Timeout:  DefaultTimeout,
		// Callback: callback,
	})
	if err != nil {
		// handle error
	}

	client.Run("ls -la")
}
func ssh1(host, cmd string) {
	config := &ssh.ClientConfig{
		// User: "hiro",
		// Auth: []ssh.AuthMethod{
		// ssh.Password("password"),
		// },
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
	}
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
