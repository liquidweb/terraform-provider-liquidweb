package liquidweb

import (
	"fmt"
	"log"
	"os"
	"strings"

	lwgoapi "github.com/liquidweb/liquidweb-go/client"
)

const lwApiURL = "https://api.liquidweb.com"
const lwApiTimeout = 15

var ErrUsernameNotSet = fmt.Errorf("liquidweb username not set")
var ErrSecretNotSet = fmt.Errorf("liquidweb secret not set")

// Config holds all of the metadata for the provider and the Liquidweb API client.
type Config struct {
	LWAPI          *lwgoapi.API
	username       string
	password       string
	URL            string
	Timeout        int
	TracingEnabled bool
}

// NewConfig accepts configuration parameters and returns a Config.
func NewConfig(username string, password string, url string, timeout int, nClient *lwgoapi.API) (*Config, error) {
	c := &Config{
		username:       username,
		password:       password,
		URL:            url,
		Timeout:        timeout,
		LWAPI:          nClient,
		TracingEnabled: len(os.Getenv("TF_LOG")) > 0,
	}

	return c, nil
}

// GetConfig handles mapping configuration between Terraform & Liquid Web's Liquidweb API configuration.
func GetConfig() (interface{}, error) {
	var username, token string

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case "LWAPI_USERNAME":
			log.Printf("read username %s\n", pair[1])
			username = pair[1]
		case "LWAPI_PASSWORD":
			token = pair[1]
		}
	}

	if username == "" {
		return nil, ErrUsernameNotSet
	}

	if token == "" {
		return nil, ErrSecretNotSet
	}

	// Initial new LW go client.
	lwAPI, err := lwgoapi.NewAPI(
		username,
		token,
		lwApiURL,
		lwApiTimeout)
	if err != nil {
		return nil, err
	}

	return NewConfig(
		username,
		token,
		lwApiURL,
		lwApiTimeout,
		lwAPI,
	)
}
