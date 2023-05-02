package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	pubkey     string
	Privkey    string
	privKeyHex string
	pin        string
	Relays     []*Relay
	follows    []*string
	Dark       bool
}

func NewConfig() *Config {
	return &Config{
		pubkey:     "",
		Privkey:    "",
		privKeyHex: "",
		pin:        "",
		Relays:     []*Relay{},
		follows:    []*string{},
		Dark:       true,
	}
}

func (c *Config) OpenConfigFile(flags int) (*os.File, error) {
	userConfigDir, _ := os.UserConfigDir()
	configDir := filepath.Join(userConfigDir, "greet")
	configPath := filepath.Join(configDir, "config.json")

	fmt.Println("User config dir", userConfigDir)
	fmt.Println("Config dir", configDir)
	fmt.Println("Config path", configPath)

	_ = os.Mkdir(configDir, 0755)
	f, err := os.OpenFile(configPath, flags, 0644)
	if err != nil {
		// Does not exist? Create
		f, err = c.OpenConfigFile(os.O_CREATE | os.O_RDWR | os.O_TRUNC)
		if err != nil {
			log.Panicln("Could not open or create the config file")
		}
		return f, err
	}
	return f, nil
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
		log.Panicln(err)
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
