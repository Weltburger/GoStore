package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Addr         string `json:"addr"`
	DBConnection string `json:"db_connection"`
	JWTKey       []byte `json:"jwt_key"`
	RedisConfig  struct {
		Addr string `json:"addr"`
		DB   int    `json:"db"`
	} `json:"redis_config"`
}

func New(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	if err := json.Unmarshal(file, config); err != nil {
		return nil, err
	}

	return config, nil
}
