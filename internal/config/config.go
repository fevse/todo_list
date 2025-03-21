package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB         DBConf
	HTTPServer HTTPConf
	TgBot      TBotConf
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
	err = viper.Unmarshal(&c)
	if err != nil {
		return Config{}, err
	}
	t, err := NewTBotConf()
	if err != nil {
		return Config{}, err
	}
	c.TgBot = t

	return
}

func NewTBotConf() (t TBotConf, err error) {
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	if err != nil {
		return TBotConf{}, err
	}
	err = viper.Unmarshal(&t)
	if err != nil {
		return TBotConf{}, err
	}
	return
}

func (c *Config) DBConnectionString() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable",
		c.DB.Type, c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)
}
