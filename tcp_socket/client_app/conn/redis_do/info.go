package redisDo

import "github.com/spf13/viper"

type RedisDo struct {
	Host string
	Port string
}

func InitRedisDo(cfg *viper.Viper) *RedisDo {
	return &RedisDo{
		Host: cfg.GetString("settings.redis.host"),
		Port: cfg.GetString("settings.redis.port"),
	}
}
