package appconfig

type MqttConfig struct {
	Config struct {
		Broker   string `yaml:"broker"`
		ClientID string `yaml:"client_id"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"config"`

	Topics map[string]MqttTopic `yaml:"topics"`
}

type MqttTopic struct {
	CommandTopic string `yaml:"command_topic"`
	StateTopic   string `yaml:"state_topic"`
}
