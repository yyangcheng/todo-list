package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ListenHost string
	ListenPort string
	MySQL      MySQLConfig
	Redis      RedisConfig
	JWT        JWTConfig
}

type MySQLConfig struct {
	Host     string
	Name     string
	Port     string
	User     string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
}

var Env *Config

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/app/config")
	viper.SetConfigType("json")
	viper.ReadInConfig()
	viper.Unmarshal(&Env)
}
