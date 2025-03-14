package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB         DBConf
	HTTPServer HTTPConf
	TgBot      TBotConf
}

type DBConf struct {
	Dir  string
	Name string
}

type HTTPConf struct {
	Host string
	Port string
}

type TBotConf struct {
	Token string
}

func NewConfig() (c Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	viper.Unmarshal(&c)
	return
}
