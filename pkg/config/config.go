package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Kafka struct {
	Mode  string `yaml:"mode"`
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Topic string `yaml:"topic"`
}

type Config struct {
	ModeLog    string `yaml:"modelog"`
	ServerPort string `yaml:"serverport"`
	DB         *DB    `yaml:"db"`
	Kafka	*Kafka 	`yaml:"kafka"`
}

func LoadConfig(fileName, fileType string) (*Config, error) {
	var cfg *Config
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(fileName)
	v.AddConfigPath("./configs")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to decode into struct, %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return cfg, nil
}
