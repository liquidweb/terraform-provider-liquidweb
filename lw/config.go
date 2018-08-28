package lw

import (
	lwapi "github.com/liquidweb/go-lwApi"
	"github.com/spf13/viper"
)

// Config holds all of the metadata for the provider and the Storm API client.
type Config struct {
	Client   *lwapi.Client
	username string
	password string
	URL      string
	Timeout  int
}

// NewConfig accepts configuration parameters and returns a Config.
func NewConfig(username string, password string, url string, timeout int, client *lwapi.Client) (*Config, error) {
	c := &Config{
		username: username,
		password: password,
		URL:      url,
		Timeout:  timeout,
		Client:   client,
	}

	return c, nil
}

// GetConfig handles mapping configuration between Terraform & Liquid Web's Storm API configuration.
func GetConfig(path string) (interface{}, error) {
	vc := viper.New()
	vc.SetConfigFile(path)
	vc.AutomaticEnv()

	verr := vc.ReadInConfig()
	if verr != nil {
		return nil, verr
	}

	client, err := lwapi.New(vc)
	if err != nil {
		return nil, err
	}

	return NewConfig(
		vc.GetString("username"),
		vc.GetString("password"),
		vc.GetString("url"),
		vc.GetInt("Timeout"),
		client,
	)
}
