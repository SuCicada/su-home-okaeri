package appconfig

type MqttConfig struct {
	Config struct {
		Broker   string `yaml:"broker"`
		ClientID string `yaml:"client_id"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"config"`

	Topics map[string]MqttTopics `yaml:"topics"`
}

type MqttTopics map[string]string

// type MqttTopic struct {
// CommandTopic string `yaml:"command_topic"`
// StateTopic   string `yaml:"state_topic"`
// }
