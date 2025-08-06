package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"port"`
	CertPath    string `mapstructure:"cert_path"`
	KeyPath     string `mapstructure:"key_path"`
	UpstreamDNS string `mapstructure:"upstream_dns"`
}

func Load() Config {
	//load config from a yaml or similar file.
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file %v", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error decoding into struct %v", err)
	}

	return cfg
}
