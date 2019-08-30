package liquidweb

import (
	"os"

	lwgoapi "github.com/liquidweb/liquidweb-go/client"
	lwapi "github.com/liquidweb/go-lwApi"
	"github.com/spf13/viper"
)

// Config holds all of the metadata for the provider and the Storm API client.
type Config struct {
	Client         *lwapi.Client
	LWAPI          *lwgoapi.API
	username       string
	password       string
	URL            string
	Timeout        int
	TracingEnabled bool
}

// NewConfig accepts configuration parameters and returns a Config.
func NewConfig(username string, password string, url string, timeout int, client *lwapi.Client, nClient *lwgoapi.API) (*Config, error) {
	c := &Config{
		username:       username,
		password:       password,
		URL:            url,
		Timeout:        timeout,
		Client:         client,
		LWAPI:          nClient,
		TracingEnabled: len(os.Getenv("TF_LOG")) > 0,
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

	username := vc.GetString("lwapi.username")
	password := vc.GetString("lwapi.password")
	url := vc.GetString("lwapi.url")
	timeout := vc.GetInt("lwapi.Timeout")

	config := lwapi.LWAPIConfig{
		Username: &username,
		Password: &password,
		Url:      url,
	}

	// Initialize original LW go client.
	client, err := lwapi.New(&config)
	if err != nil {
		return nil, err
	}

	// Initial new LW go client.
	lwAPI, err := lwgoapi.NewAPI(username, password, url, timeout)
	if err != nil {
		return nil, err
	}

	return NewConfig(
		username,
		password,
		url,
		timeout,
		client,
		lwAPI,
	)
}
