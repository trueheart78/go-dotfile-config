package dotfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var configFilename = ".go-call-me.json"

// Config is the object that holds the json values
type Config struct {
	RedisURL      string `json:"redis_url"`
	RedisPassword string `json:"redis_password"`
	RedisChannels struct {
		Emergency   string `json:"emergency"`
		NonEmergent string `json:"nonemergent"`
	} `json:"redis_channels"`
	Path   string
	Loaded bool
}

// NewConfig attempts to load an existing configuration
func NewConfig() (c Config, err error) {
	fullConfig := configPath(configFilename)
	if _, err = os.Stat(fullConfig); os.IsNotExist(err) {
		err = fmt.Errorf("Please setup your config file [%v]", fullConfig)
		return
	}
	_, err = c.load(fullConfig)
	if err != nil {
		return
	}
	if !c.Valid() {
		err = fmt.Errorf("please validate the %v file. See README for details", fullConfig)
	}
	return
}

func configPath(filename string) string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, filename)
}

func configExists(configFilename string) bool {
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		return false
	}
	return true
}

// Valid returns whether the config is valid
func (c Config) Valid() (ok bool) {
	ok, _ = c.validate()
	return
}

func (c Config) validate() (ok bool, err error) {
	if !c.Loaded {
		err = errors.New("the config has yet to be loaded")
		return
	}
	if c.RedisURL == "" || c.RedisPassword == "" || c.RedisChannels.Emergency == "" || c.RedisChannels.NonEmergent == "" {
		err = errors.New("Incomplete configuration")
		return
	}
	ok = true
	return
}

func (c *Config) load(configFilename string) (ok bool, err error) {
	var raw []byte
	raw, err = ioutil.ReadFile(configFilename)
	if err != nil {
		c.Path = configFilename
		return
	}
	json.Unmarshal(raw, c)
	ok = true
	c.Path = configFilename
	if configExists(configFilename) {
		c.Loaded = true
	}
	return
}
