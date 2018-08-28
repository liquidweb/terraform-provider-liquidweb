package lw

import (
	lwgo "git.liquidweb.com/masre/liquidweb-go/client"
	lwapi "github.com/liquidweb/go-lwApi"
	"github.com/spf13/viper"
)

// Config holds all of the metadata for the provider and the Storm API client.
type Config struct {
	Client   *lwapi.Client
	NClient  *lwgo.Client
	username string
	password string
	URL      string
	Timeout  int
}

// NewConfig accepts configuration parameters and returns a Config.
func NewConfig(username string, password string, url string, timeout int, client *lwapi.Client, nClient *lwgo.Client) (*Config, error) {
	c := &Config{
		username: username,
		password: password,
		URL:      url,
		Timeout:  timeout,
		Client:   client,
		NClient:  nClient,
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

	username := vc.GetString("username")
	password := vc.GetString("password")
	url := vc.GetString("url")
	timeout := vc.GetInt("Timeout")

	// Initialize original LW go client.
	client, err := lwapi.New(vc)
	if err != nil {
		return nil, err
	}

	// Initial new LW go client.
	nConfig, err := lwgo.NewConfig(username, password, url, timeout, true)
	if err != nil {
		return nil, err
	}
	nClient := lwgo.NewClient(nConfig)

	return NewConfig(
		username,
		password,
		url,
		timeout,
		client,
		nClient,
	)
}
