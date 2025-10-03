package mqttpkg

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/util"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttClient mqtt.Client
)

func Init() {
	log.Println("🚀 启动 Linux Display MQTT 桥接程序")

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
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(onConnect)
	opts.SetConnectionLostHandler(onConnectionLost)

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ MQTT 连接失败: %v", token.Error())
	}

	log.Println("✅ MQTT 连接成功")

}

func Close() {
	mqttClient.Disconnect(250)
}

func RegisterRoute(topic string, callback mqtt.MessageHandler) {
	routes <- MqttRoute{Topic: topic, Qos: 1, Callback: callback}
}

type MqttRoute struct {
	Topic    string
	Qos      byte
	Callback mqtt.MessageHandler
}

var routes = make(chan MqttRoute, 10)

func onConnect(client mqtt.Client) {
	log.Println("🔗 MQTT 客户端已连接")

	for route := range routes {
		// 订阅命令主题
		if token := client.Subscribe(route.Topic, 1, route.Callback); token.Wait() && token.Error() != nil {
			log.Printf("❌ 订阅 %s 失败: %v", route.Topic, token.Error())
		} else {
			log.Printf("📡 已订阅: %s", route.Topic)
		}
	}

}

func onConnectionLost(client mqtt.Client, err error) {
	log.Printf("❌ MQTT 连接丢失: %v", err)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ MQTT 连接失败: %v", token.Error())
	}
	log.Println("✅ MQTT 连接成功")
}

func messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("📨 收到消息: %s = %s", msg.Topic(), msg.Payload())
}
func Publish(topic string, payload interface{}) {
	payloadBytes := util.Conv.ToBytes(payload)

	log.Printf("📨 发布消息: %s = %s", topic, string(payloadBytes))
	if token := mqttClient.Publish(topic, 1, true, payloadBytes); token.Wait() && token.Error() != nil {
		log.Printf("❌ 发布消息失败: %v", token.Error())
	}
}

// func handlePowerCommand(client mqtt.Client, msg mqtt.Message) {
// 	command := string(msg.Payload())
// 	log.Printf("🔌 收到开关命令: %s", command)

// 	switch command {
// 	case "ON":
// 		if callHTTPAPI("turn_on", "POST", nil) {
// 			currentState.IsOn = "ON"
// 			publishState()
// 		}
// 	case "OFF":
// 		if callHTTPAPI("turn_off", "POST", nil) {
// 			currentState.IsOn = "OFF"
// 			currentState.Brightness = 0
// 			publishState()
// 		}
// 	default:
// 		log.Printf("⚠️ 未知开关命令: %s", command)
// 	}
// }

// func handleBrightnessCommand(client mqtt.Client, msg mqtt.Message) {
// 	brightnessStr := string(msg.Payload())
// 	brightness, err := strconv.Atoi(brightnessStr)
// 	if err != nil {
// 		log.Printf("❌ 亮度值解析失败: %s", brightnessStr)
// 		return
// 	}

// 	log.Printf("💡 收到亮度命令: %d", brightness)

// 	// 调用 HTTP API 设置亮度
// 	data := APIRequest{
// 		Brightness: brightness,
// 		Value:      brightness, // 有些 API 用 value 字段
// 	}

// 	if callHTTPAPI("set_brightness", "POST", data) { // 修改为你的亮度设置端点
// 		currentState.Brightness = brightness
// 		if brightness > 0 {
// 			currentState.IsOn = "ON"
// 		}
// 		// publishState()
// 	}
// }

// func callHTTPAPI(endpoint, method string, data interface{}) bool {
// 	url := fmt.Sprintf("%s/%s", HTTPAPIBase, endpoint) // 修改为你的 API 路径格式

// 	var body io.Reader
// 	if data != nil {
// 		jsonData, err := json.Marshal(data)
// 		if err != nil {
// 			log.Printf("❌ JSON 序列化失败: %v", err)
// 			return false
// 		}
// 		body = bytes.NewBuffer(jsonData)
// 	}

// 	client := &http.Client{Timeout: HTTPTimeout}
// 	req, err := http.NewRequest(method, url, body)
// 	if err != nil {
// 		log.Printf("❌ 创建 HTTP 请求失败: %v", err)
// 		return false
// 	}

// 	if data != nil {
// 		req.Header.Set("Content-Type", "application/json")
// 	}

// 	if DebugMode {
// 		log.Printf("🌐 调用 API: %s %s", method, url)
// 		if data != nil {
// 			jsonData, _ := json.Marshal(data)
// 			log.Printf("📤 请求数据: %s", jsonData)
// 		}
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("❌ HTTP 请求失败: %v", err)
// 		return false
// 	}
// 	defer resp.Body.Close()

// 	respBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("❌ 读取响应失败: %v", err)
// 		return false
// 	}

// 	if DebugMode {
// 		log.Printf("📥 API 响应: %s", respBody)
// 	}

// 	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
// 		log.Printf("✅ API 调用成功: %s", endpoint)
// 		return true
// 	} else {
// 		log.Printf("❌ API 调用失败: %d %s", resp.StatusCode, respBody)
// 		return false
// 	}
// }

// // func publishState() {
// // 	if mqttClient == nil || !mqttClient.IsConnected() {
// // 		log.Println("⚠️ MQTT 客户端未连接，跳过状态发布")
// // 		return
// // 	}

// // 	// 发布开关状态
// // 	if token := mqttClient.Publish(StateTopic, 1, true, currentState.IsOn); token.Wait() && token.Error() != nil {
// // 		log.Printf("❌ 发布开关状态失败: %v", token.Error())
// // 	}

// // 	// 发布亮度状态
// // 	brightnessStr := strconv.Itoa(currentState.Brightness)
// // 	if token := mqttClient.Publish(BrightnessStateTopic, 1, true, brightnessStr); token.Wait() && token.Error() != nil {
// // 		log.Printf("❌ 发布亮度状态失败: %v", token.Error())
// // 	}

// // 	log.Printf("✅ 状态已发布: 开关=%s, 亮度=%d", currentState.IsOn, currentState.Brightness)

// // 	if DebugMode {
// // 		log.Printf("📡 发布到 %s: %s", StateTopic, currentState.IsOn)
// // 		log.Printf("📡 发布到 %s: %s", BrightnessStateTopic, brightnessStr)
// // 	}
// // }

// // 可选：定期从 HTTP API 同步状态
// func getDeviceStatus() {
// 	resp, err := http.Get(fmt.Sprintf("%s/status", HTTPAPIBase)) // 修改为你的状态查询端点
// 	if err != nil {
// 		if DebugMode {
// 			log.Printf("⚠️ 获取设备状态失败: %v", err)
// 		}
// 		return
// 	}
// 	defer resp.Body.Close()

// 	var apiResp APIResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
// 		if DebugMode {
// 			log.Printf("⚠️ 解析状态响应失败: %v", err)
// 		}
// 		return
// 	}

// 	// 更新状态（根据你的 API 响应格式调整）
// 	if apiResp.Power {
// 		currentState.IsOn = "ON"
// 	} else {
// 		currentState.IsOn = "OFF"
// 	}
// 	currentState.Brightness = apiResp.Brightness

// 	publishState()
// }
