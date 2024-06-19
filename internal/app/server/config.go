package server

import (
	"encoding/json"
	"errors"
	"os"
)

// Config represents the configuration structure for the server.
type Config struct {
	Port     string `json:"port"`
	DbPath   string `json:"db_path"`
	DbDriver string `json:"db_driver"`
}

// NewConfig creates a new Config instance with default values.
func NewConfig() Config {
	return Config{}
}

// ReadConfig reads the configuration from a JSON file and populates the Config struct.
func (c *Config) ReadConfig() error {
	// Read JSON data from the config file
	jsonData, err := os.ReadFile("./configs/config.json")
	if err != nil {
		return errors.Join(errors.New("error reading config file"), err)
	}

	// Unmarshal JSON data into the Config struct
	err = json.Unmarshal(jsonData, c)
	if err != nil {
		return errors.Join(errors.New("error unmarshalling JSON"), err)
	}
	return nil
}
