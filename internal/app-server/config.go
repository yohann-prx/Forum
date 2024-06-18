package server

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port     string `json:"port"`
	DbPath   string `json:"db_path"`
	DbDriver string `json:"db_driver"`
}

func NewConfig() Config {
	return Config{}
}
func (c *Config) ReadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}

	return nil
}
