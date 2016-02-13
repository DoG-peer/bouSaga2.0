package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config carries information
type Config struct {
	URL   string
	Res   int
	Cache string
	Jtalk Jtalk
	Voice []VoiceConfig
}

func (c Config) String() string {
	s, _ := json.MarshalIndent(c, "", "  ")
	return string(s)
}

// Save saves config
func (c *Config) Save(fname string) error {
	return ioutil.WriteFile(fname, []byte(c.String()), os.ModePerm)
}

// Load loads config
func (c *Config) Load(fname string) error {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return err
	}
	return nil
}
