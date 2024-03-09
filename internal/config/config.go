package config

import (
	"errors"
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	StoragePath string        `yaml:"storage_path"`
	TokenTTL    time.Duration `yaml:"token_ttl"`
	HTTP
}

type HTTP struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeOut  time.Duration `yaml:"read_timeout"`
	WriteTimeOut time.Duration `yaml:"write_timeout"`
}

// Load load config, panic if has error
func Load() *Config {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.Parse()

	if cfgPath == "" {
		cfgPath = "./config/local.yaml"
	}

	if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		panic("config file is not exists in: " + cfgPath)
	}
	var cfg Config
	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
