package dotfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConfigFilename = fixturePath("valid")
var missingAuthConfigFilename = fixturePath("missing_auth")
var missingURLConfigFilename = fixturePath("missing_url")
var missingChannelsConfigFilename = fixturePath("missing_channels")
var missingConfigFilename = fixturePath("missing")
var emptyConfigFilename = fixturePath("empty")

func TestConfigPath(t *testing.T) {
	fullConfigPath := configPath(".go-call-me.json")

	assert.True(t, strings.HasSuffix(fullConfigPath, configFilename))
	assert.False(t, strings.HasPrefix(fullConfigPath, configFilename))
}

func TestConfigExists(t *testing.T) {
	assert.True(t, configExists(validConfigFilename))
	assert.True(t, configExists(emptyConfigFilename))
	assert.False(t, configExists(missingConfigFilename))
}

func TestConfigLoad(t *testing.T) {
	assert := assert.New(t)

	// valid config
	c := Config{}
	c.load(validConfigFilename)

	assert.Equal("https://url.to.redis:1234", c.RedisURL)
	assert.Equal("password", c.RedisPassword)
	assert.Equal("emergency", c.RedisChannels.Emergency)
	assert.Equal("nonemergent", c.RedisChannels.NonEmergent)
	assert.Equal(validConfigFilename, c.Path)
	assert.True(c.Loaded)

	// empty config
	c = Config{}
	c.load(emptyConfigFilename)

	assert.Equal("", c.RedisURL)
	assert.Equal("", c.RedisPassword)
	assert.Equal("", c.RedisChannels.Emergency)
	assert.Equal("", c.RedisChannels.NonEmergent)
	assert.Equal(emptyConfigFilename, c.Path)
	assert.True(c.Loaded)

	// missing config
	c = Config{}
	c.load(missingConfigFilename)

	assert.Equal("", c.RedisURL)
	assert.Equal("", c.RedisPassword)
	assert.Equal("", c.RedisChannels.Emergency)
	assert.Equal("", c.RedisChannels.NonEmergent)
	assert.Equal(missingConfigFilename, c.Path)
	assert.False(c.Loaded)
}

func TestConfigValidate(t *testing.T) {
	assert := assert.New(t)

	c := Config{}
	c.load(validConfigFilename)
	ok, err := c.validate()
	assert.Nil(err)
	assert.True(c.Valid())
	assert.True(ok)
	assert.Nil(err)

	c = Config{}
	c.load(missingURLConfigFilename)
	ok, err = c.validate()
	assert.False(c.Valid())
	assert.False(ok)
	assert.NotNil(err)
	assert.Equal("Incomplete configuration", err.Error())

	c = Config{}
	c.load(missingAuthConfigFilename)
	ok, err = c.validate()
	assert.False(c.Valid())
	assert.False(ok)
	assert.NotNil(err)
	assert.Equal("Incomplete configuration", err.Error())

	c = Config{}
	c.load(missingChannelsConfigFilename)
	ok, err = c.validate()
	assert.False(c.Valid())
	assert.False(ok)
	assert.Equal("Incomplete configuration", err.Error())

	c = Config{}
	c.load(emptyConfigFilename)
	ok, err = c.validate()
	assert.False(c.Valid())
	assert.False(ok)
	assert.NotNil(err)
}

func fixturePath(filename string) string {
	workingDir, _ := os.Getwd()
	return fmt.Sprintf("%v.json", filepath.Join(workingDir, "fixtures", filename))
}
