package cmdutil

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RabbitmqHost     string `yaml:"rabbitmq_host"`
	RabbitmqPort     uint   `yaml:"rabbitmq_port"`
	RabbitmqUser     string `yaml:"rabbitmq_user"`
	RabbitmqPassword string `yaml:"rabbitmq_password"`
	SlackWebHookUrl  string `yaml:"slack_webhook_url"`
}

func ReadConfigFile(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Uri() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitmqUser, c.RabbitmqPassword, c.RabbitmqHost, c.RabbitmqPort)
}
