package main

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type Config struct {
	pubkey        string
	Privkey       string
	privKeyHex    string
	pin           string
	Relays        []*Relay
	follows       []*string
	Dark          bool
	userConfigDir string
	configDir     string
	configPath    string
}

func NewConfig() *Config {
	userConfigDir, _ := os.UserConfigDir()
	configDir := filepath.Join(userConfigDir, "greet")
	configPath := filepath.Join(configDir, "config.json")

	log.Debug().Msgf("User config dir %s", userConfigDir)
	log.Debug().Msgf("Config dir %s", configDir)
	log.Debug().Msgf("Config path %s", configPath)

	return &Config{
		pubkey:        "",
		Privkey:       "",
		privKeyHex:    "",
		pin:           "",
		Relays:        []*Relay{},
		follows:       []*string{},
		Dark:          true,
		userConfigDir: userConfigDir,
		configDir:     configDir,
		configPath:    configPath,
	}
}

func (c *Config) OpenConfigFile(flags int) (*os.File, error) {
	_ = os.Mkdir(c.configDir, 0755)
	return openFile(c.configPath, flags)
}

func (c *Config) Load() error {

	f, err := c.OpenConfigFile(os.O_RDONLY)
	if err != nil {
		return err
	}
	defer f.Close()

	fileinfo, err := f.Stat()
	if err != nil {
		return err
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	bytesRead, err := f.Read(buffer)
	if err != nil {
		return err
	}

	if bytesRead == 0 {
		var j []byte
		f.Close()
		f, err = c.OpenConfigFile(os.O_CREATE | os.O_RDWR | os.O_TRUNC)
		if err != nil {
			return err
		}
		j, err = json.Marshal(c)
		if err != nil {
			return err
		}
		f.Write(j)
		return nil
	}

	err = json.Unmarshal(buffer, &c)
	if err != nil {
		log.Err(err)
	}

	return nil
}

func (c *Config) Save() error {
	f, err := c.OpenConfigFile(os.O_CREATE | os.O_RDWR | os.O_TRUNC)
	if err != nil {
		return err
	}
	defer f.Close()
	configOutput, err := PrettyStruct(c)
	if err != nil {
		return err
	}

	f.Write([]byte(configOutput + "\n"))

	return nil
}
