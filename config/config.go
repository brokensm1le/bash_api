package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
}

type ServerConfig struct {
	AppVersion string `json:"appVersion"`
	Host       string `json:"host" validate:"required"`
	Port       string `json:"port" validate:"required"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"-"`
	DBName   string `json:"DBName"`
	SSLMode  string `json:"sslMode"`
	PgDriver string `json:"pgDriver"`
}

func LoadConfig() (*viper.Viper, error) {

	viperInstance := viper.New()

	viperInstance.AddConfigPath("./config")
	viperInstance.SetConfigName("config")
	viperInstance.SetConfigType("yml")

	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return viperInstance, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
		return nil, err
	}
	return &c, nil
}
