package liquidweb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kkyr/fig"
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
func GetConfig(resourceData *schema.ResourceData) (interface{}, error) {
	configPaths := []string{
		"./",
		"/",
	}
	if homedir, err := os.UserHomeDir(); err != nil {
		configPaths = append(configPaths, homedir)
	}

	configName := "lwapi.toml"

	if optConfigFile, isSet := resourceData.GetOk("config_path"); isSet {
		if val, ok := optConfigFile.(string); ok {
			log.Printf("included config file: %s\n", val)
			configPaths = append([]string{filepath.Dir(val)}, configPaths...)
			configName = filepath.Base(val)
		}
	}

	var cfg struct {
		lwApiCfg struct {
			Username string `fig:"username"`
			Password string `fig:"password"`
		} `fig:"lwapi"`
	}

	if err := fig.Load(&cfg,
		fig.File(configName),
		fig.Dirs(configPaths...),
	); err != nil && !errors.Is(err, fig.ErrFileNotFound) {
		return nil, err
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case "LWAPI_USERNAME":
			cfg.lwApiCfg.Username = pair[1]
		case "LWAPI_PASSWORD":
			cfg.lwApiCfg.Password = pair[1]
		}
	}

	if cfg.lwApiCfg.Username == "" {
		return nil, ErrUsernameNotSet
	}

	if cfg.lwApiCfg.Password == "" {
		return nil, ErrSecretNotSet
	}

	// Initial new LW go client.
	lwAPI, err := lwgoapi.NewAPI(
		cfg.lwApiCfg.Username,
		cfg.lwApiCfg.Password,
		lwApiURL,
		lwApiTimeout)
	if err != nil {
		return nil, err
	}

	return NewConfig(
		cfg.lwApiCfg.Username,
		cfg.lwApiCfg.Password,
		lwApiURL,
		lwApiTimeout,
		lwAPI,
	)
}
