package appConfig

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	//====BRIDGE====//
	Bridge string `envconfig:"BRIDGE"`

	//====FIRST_APP====//
	Port string `envconfig:"FIRST_APP_PORT"`

	//====SECOND_APP====//
	SecondPort string `envconfig:"SECOND_APP_PORT"`
	SecondLink string `envconfig:"SECOND_APP_LINK"`

	//====MYSQL====//
	SQLConnectLink string `envconfig:"SQL_CONNECT_LINK"`

	//====RMQ====//
	RMQConnectLink string `envconfig:"RMQ_CONNECT_LINK"`
}

func (c *Config) Set() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	*c = config
}

// Get...
func Get() Config {
	var config Config
	config.Set()
	return config
}
