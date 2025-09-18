package test

import (
	"testing"
)

//func TestSleep(m *testing.M) {
//	godotenv.Load("../.env")
//	// 加载YAML配置
//	cfg.LoadConfig("../config.yaml")
//	os.Exit(m.Run())
//}

func TestSleep(t *testing.T) {
}

func TestSsh(t *testing.T) {
	// 使用TOML配置获取Linux SSH配置
	//linuxConfig := util.GetSSHConfig("linux")
	//
	//client, err := goph.NewConn(&goph.Config{
	//	User: linuxConfig.User,
	//	Addr: linuxConfig.Host,
	//	Port: uint(linuxConfig.Port),
	//	// Auth: auth,
	//	// Timeout:  DefaultTimeout,
	//	Callback: ssh.InsecureIgnoreHostKey(),
	//})
	//if err != nil {
	//	assert.Nil(t, err)
	//}
	//
	//output, err := client.Run("ls -la")
	//if err != nil {
	//	assert.Nil(t, err)
	//}
	//fmt.Println(string(output))
}
