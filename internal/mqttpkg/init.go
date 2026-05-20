package mqttpkg

import (
	"fmt"
	"log"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/util"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttClient mqtt.Client
	routes     []MqttRoute
	routesLock sync.RWMutex
)

func Init() {
	logger.Info("🚀 启动 Linux Display MQTT 桥接程序")

	// 创建 MQTT 客户端
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.GetConfig().MQTT.Config.Broker)
	opts.SetClientID(cfg.GetConfig().MQTT.Config.ClientID)

	var username, password = cfg.GetConfig().MQTT.Config.Username, cfg.GetConfig().MQTT.Config.Password
	if username != "" && password != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetConnectTimeout(5 * time.Second)

	// 开启自动重连
	opts.SetAutoReconnect(true)
	// 设置最大重连间隔为5秒
	opts.SetMaxReconnectInterval(5 * time.Second)

	opts.SetOnConnectHandler(onConnect)
	opts.SetConnectionLostHandler(onConnectionLost)

	mqttClient = mqtt.NewClient(opts)

	// 初始连接循环
	for {
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			logger.Error("❌ MQTT 连接失败: %v", token.Error())
			logger.Info("🔄 5秒后重试...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	logger.Info("✅ MQTT 连接成功")
}

func Close() {
	if mqttClient != nil {
		mqttClient.Disconnect(250)
	}
}

func RegisterRoute(topic string, callback mqtt.MessageHandler) {
	route := MqttRoute{Topic: topic, Qos: 1, Callback: callback}

	routesLock.Lock()
	routes = append(routes, route)
	routesLock.Unlock()

	// 如果已经连接，立即订阅
	if mqttClient != nil && mqttClient.IsConnected() {
		if token := mqttClient.Subscribe(route.Topic, route.Qos, route.Callback); token.Wait() && token.Error() != nil {
			log.Printf("❌ 订阅 %s 失败: %v", topic, token.Error())
		} else {
			log.Printf("📡 已订阅: %s", topic)
		}
	}
}

type MqttRoute struct {
	Topic    string
	Qos      byte
	Callback mqtt.MessageHandler
}

func onConnect(client mqtt.Client) {
	logger.Info("🔗 MQTT 客户端已连接")

	routesLock.RLock()
	defer routesLock.RUnlock()

	for _, route := range routes {
		if token := client.Subscribe(route.Topic, route.Qos, route.Callback); token.Wait() && token.Error() != nil {
			logger.Error("❌ 订阅 %s 失败: ", route.Topic, token.Error())
		} else {
			logger.Info("📡 已订阅: ", route.Topic)
		}
	}
}

func onConnectionLost(client mqtt.Client, err error) {
	logger.Error("❌ MQTT 连接丢失: ", err)
	// 自动重连已开启，无需手动重连
}

func messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	logger.Info(fmt.Sprintf("📨 收到消息: %s = %s", msg.Topic(), msg.Payload()))
}

func Publish(topic string, payload interface{}) {
	payloadBytes := util.Conv.ToBytes(payload)
	logPayload := string(payloadBytes)
	if len(logPayload) > 1000 {
		logPayload = logPayload[:1000]
	}
	logger.Info(fmt.Sprintf("📨 发布消息: %s = %s", topic, logPayload))

	if mqttClient == nil || !mqttClient.IsConnected() {
		logger.Error("❌ 发布消息失败: MQTT 未连接")
		return
	}

	if token := mqttClient.Publish(topic, 1, true, payloadBytes); token.Wait() && token.Error() != nil {
		logger.Error("❌ 发布消息失败: ", token.Error())
	}
}
