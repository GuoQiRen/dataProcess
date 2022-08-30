package config

import "github.com/spf13/viper"

type NginxInfo struct {
	Head        string
	Host        string
	Port        string
	SaveContext string
	DescContext string
}

func InitNginxInfo(cfg *viper.Viper) *NginxInfo {
	return &NginxInfo{
		Head:        cfg.GetString("head"),
		Host:        cfg.GetString("host"),
		Port:        cfg.GetString("port"),
		SaveContext: cfg.GetString("saveContext"),
		DescContext: cfg.GetString("descContext"),
	}
}
