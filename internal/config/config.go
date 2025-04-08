package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB         DBConf
	HTTPServer HTTPConf
}

type DBConf struct {
	Type     string
	Dir      string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type HTTPConf struct {
	Host string
	Port string
}

func NewConfig() (c Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		return Config{}, err
	}

	return
}

func (c *Config) DBConnectionString() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable",
		c.DB.Type, c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)
}
