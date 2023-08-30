package liquidweb

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	lwgoapi "github.com/liquidweb/liquidweb-go/client"
)

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
func GetConfig(path string) (interface{}, error) {
	log.Printf("config file: %s\n", path)

	type lwapiConfig struct {
		Username string `toml:"username"`
		Password string `toml:"password"`
		URL      string `toml:"url"`
		Timeout  int    `toml:"timeout"`
	}

	conf := struct {
		LWAPI lwapiConfig `toml:"lwapi"`
	}{}

	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	log.Printf("conf file contained: %s\n\n", body)

	metadata, err := toml.Decode(string(body), &conf)
	if err != nil {
		return nil, err
	}
	_ = metadata

	// config := lwapi.LWAPIConfig{
	// 	Username: &conf.LWAPI.Username,
	// 	Password: &conf.LWAPI.Password,
	// 	Url:      conf.LWAPI.URL,
	// }

	// vc := viper.New()
	// vc.SetConfigFile(path)
	// vc.AutomaticEnv()
	// vc.SetConfigType("toml")

	// if err = vc.ReadConfig(bytes.NewBuffer(body)); err != nil {
	// 	return nil, err
	// }

	// log.Printf("configuration keys %#v exist\n", vc.AllKeys())
	// for _, k := range vc.AllKeys() {
	// 	log.Printf("setting %s is set to %#v\n", k, vc.Get(k))
	// }

	// username := vc.GetString("lwapi.username")
	// password := vc.GetString("lwapi.password")
	// url := vc.GetString("lwapi.url")
	// timeout := vc.GetInt("lwapi.Timeout")

	// config := lwapi.LWAPIConfig{
	// 	Username: &username,
	// 	Password: &password,
	// 	Url:      url,
	// }

	// Initialize original LW go client.
	// client, err := lwapi.New(&config)
	// if err != nil {
	// 	return nil, err
	// }

	// Initial new LW go client.
	lwAPI, err := lwgoapi.NewAPI(conf.LWAPI.Username, conf.LWAPI.Password, conf.LWAPI.URL, conf.LWAPI.Timeout)
	if err != nil {
		return nil, err
	}

	return NewConfig(
		conf.LWAPI.Username,
		conf.LWAPI.Password,
		conf.LWAPI.URL,
		conf.LWAPI.Timeout,
		lwAPI,
	)
}
