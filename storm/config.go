package storm

import (
	lwapi "github.com/liquidweb/go-lwApi"
	"github.com/spf13/viper"
)

// Config holds all of the metadata for the provider and the Storm API client.
type Config struct {
	username string
	password string
	URL      string
	Timeout  int
}

// Client returns a new API client that is used by various resources.
func (c *Config) Client() (interface{}, error) {
	vc := viper.New()
	vc.Set("username", c.username)
	vc.Set("password", c.password)
	client, err := lwapi.New(vc)

	return &client, err
}

// GetConfig handles mapping configuration between Terraform & Liquid Web's Storm API configuration.
func GetConfig(path string) (*Config, error) {
	vconfig := viper.New()
	vconfig.SetConfigFile(path)
	vconfig.AutomaticEnv()

	verr := vconfig.ReadInConfig()
	if verr != nil {
		return nil, verr
	}

	config := &Config{
		username: vconfig.GetString("username"),
		password: vconfig.GetString("password"),
		URL:      vconfig.GetString("url"),
		Timeout:  vconfig.GetInt("Timeout"),
	}

	return config, nil
}
