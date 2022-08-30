package config

import "github.com/spf13/viper"

type Cache struct {
	CacheType string
	Host      string
	Port      string
}

func InitCache(cfg *viper.Viper) *Cache {
	return &Cache{
		Host: cfg.GetString("host"),
		Port: cfg.GetString("port"),
	}
}
